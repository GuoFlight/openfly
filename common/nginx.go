package common

import (
	"encoding/json"
	"fmt"
	"github.com/GuoFlight/gerror"
	"github.com/sirupsen/logrus"
	"openfly/conf"
	"openfly/logger"
	"openfly/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Nginx struct {
	Lock sync.Mutex
}

var GNginx Nginx

// todo
func (n *Nginx) CheckConfigL4(l4 []NginxConfL4) *gerror.Gerr {
	return nil
}
func (n *Nginx) GenFilePathL4(l4 NginxConfL4) string {
	return path.Join(conf.PathData, strconv.Itoa(l4.Listen)+conf.NgFileExtension)
}
func (n *Nginx) GenConfigL4(l4 NginxConfL4) (string, *gerror.Gerr) {
	// 校验参数
	gerr := n.CheckConfigL4([]NginxConfL4{l4})
	if gerr != nil {
		return "", gerr
	}
	// 开始生成server块
	var confServer []string
	// 生成upstream
	confUpstream := n.genConfigL4Upstream(l4.Upstream, l4.Listen)
	// 生成server语句
	confServer = append(confServer, fmt.Sprintf("listen %d;", l4.Listen))
	// 生成proxy_pass语句
	confServer = append(confServer, fmt.Sprintf("proxy_pass %d;", l4.Listen))
	// 生成include语句
	if len(l4.IncludeFiles) > 0 {
		var confIncludeFiles []string
		for _, includeFile := range l4.IncludeFiles {
			confIncludeFiles = append(confIncludeFiles, fmt.Sprintf("include %s;", includeFile))
		}
		confServer = append(confServer, strings.Join(confIncludeFiles, "\n\t"))
	}
	// 生成白名单
	if len(l4.WhiteList) > 0 {
		var confWhiteList []string
		for _, whiteItem := range l4.WhiteList {
			confWhiteList = append(confWhiteList, fmt.Sprintf("%s %s;", whiteItem.Type, whiteItem.Target))
		}
		confServer = append(confServer, strings.Join(confWhiteList, "\n\t"))
	}
	return fmt.Sprintf("%s\nserver{\n\t%s\n}", confUpstream, strings.Join(confServer, "\n\t")), nil

}
func (n *Nginx) genConfigL4Upstream(upstream Upstream, port int) string {
	var confUpstream []string
	// 生成upstream中的server语句
	var hosts []string
	for _, host := range upstream.Hosts {
		if host.Weight <= 0 {
			host.Weight = conf.GConf.Nginx.WeightDefault
		}
		if host.MaxFails <= 0 {
			host.MaxFails = conf.GConf.Nginx.MaxFailsDefault
		}
		if host.FailTimeoutSecond <= 0 {
			host.FailTimeoutSecond = conf.GConf.Nginx.FailTimeoutDefault
		}
		server := fmt.Sprintf(
			"server %s:%d weight=%d max_fails=%d fail_timeout=%ds",
			host.Ip, host.Port, host.Weight, host.MaxFails, host.FailTimeoutSecond,
		)
		if host.IsBackup {
			server = server + " backup;"
		} else {
			server = server + ";"
		}
		hosts = append(hosts, server)
	}
	confUpstream = append(confUpstream, strings.Join(hosts, "\n\t"))
	// 判断是否哈希
	if upstream.IsHash {
		if upstream.HashField == "" {
			upstream.HashField = "remote_addr"
		}
		confUpstream = append(confUpstream, fmt.Sprintf("hash $%s consistent;", upstream.HashField))
	}
	// 主动健康检查
	if upstream.Interval > 0 {
		confUpstream = append(confUpstream, fmt.Sprintf(
			"check interval=%d rise=%d fall=%d timeout=%d type=tcp;",
			upstream.Interval, upstream.Rise, upstream.Fall, upstream.Timeout))
	}
	// 返回结果
	return fmt.Sprintf("upstream %d {\n\t%s\n}", port, strings.Join(confUpstream, "\n\t"))
}

func (n *Nginx) writeFile(l4 NginxConfL4) *gerror.Gerr {
	nginxConfL4, gerr := n.GenConfigL4(l4)
	if gerr != nil {
		return gerr
	}
	targetFile := n.GenFilePathL4(l4)
	err := os.WriteFile(targetFile, []byte(nginxConfL4), 0660)
	if err != nil {
		return gerror.NewErr(err.Error())
	}
	return nil
}
func (n *Nginx) writeFileAndReload(l4 NginxConfL4) *gerror.Gerr {
	gerr := n.writeFile(l4)
	if gerr != nil {
		return gerr
	}
	_, gerr = n.testAndReload()
	return gerr
}
func (n *Nginx) delFileAndReload(l4 NginxConfL4) *gerror.Gerr {
	targetFile := n.GenFilePathL4(l4)
	err := utils.DelFile(targetFile)
	if err != nil {
		return gerror.NewErr(err.Error())
	}
	_, gerr := n.testAndReload()
	return gerr
}
func (n *Nginx) test() (string, *gerror.Gerr) {
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), gerror.NewErr(string(output) + " " + err.Error())
	}
	return string(output), nil
}
func (n *Nginx) reload() (string, *gerror.Gerr) {
	cmd := exec.Command("nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), gerror.NewErr(err.Error())
	}
	return string(output), nil
}
func (n *Nginx) testAndReload() (string, *gerror.Gerr) {
	output, gerr := n.test()
	if gerr != nil {
		return output, gerr
	}
	return n.reload()
}

func (n *Nginx) WriteFileAndReload(l4 NginxConfL4) *gerror.Gerr {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	return n.writeFileAndReload(l4)
}
func (n *Nginx) DelFileAndReload(l4 NginxConfL4) *gerror.Gerr {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	return n.delFileAndReload(l4)
}

// Reset 根据etcd重置nginx配置
func (n *Nginx) Reset() *gerror.Gerr {
	// 从etcd中得到所有配置
	AllConfig, gerr := n.GetAll()
	if gerr != nil {
		return gerr
	}
	// 备份当前nginx中的配置
	n.Lock.Lock()
	defer n.Lock.Unlock()
	_, gerr = BackupFile(conf.PathData)
	if gerr != nil {
		return gerr
	}
	// 删除当前nginx中的配置
	err := filepath.Walk(conf.PathData, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if conf.PathData == path {
			return nil
		}
		err = utils.DelFile(path)
		if err != nil {
			return gerr
		}
		return nil
	})
	if err != nil {
		logger.GLogger.Fatal(err)
	}
	// 写入配置
	for _, conf := range AllConfig {
		gerr := n.writeFileAndReload(conf)
		if gerr != nil {
			logger.GLogger.Fatal(gerr)
		}
	}
	return nil
}
func (n *Nginx) GetAll() ([]NginxConfL4, *gerror.Gerr) {
	kvs, gerr := GEtcd.GetAllL4()
	if gerr != nil {
		return nil, gerr
	}
	var l4s []NginxConfL4
	for _, kv := range kvs {
		var l4 NginxConfL4
		err := json.Unmarshal(kv.Value, &l4)
		if err != nil {
			logger.GLogger.WithFields(logrus.Fields{"key": string(kv.Key), "value": string(kv.Value)}).Error("解析json失败")
			continue
		}
		l4s = append(l4s, l4)
	}
	return l4s, nil
}
