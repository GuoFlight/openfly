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
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Nginx struct {
	Lock sync.Mutex
}

var GNginx Nginx

func (n *Nginx) IsValidWhiteList(wls []WhiteListItem) *gerror.Gerr {
	for _, wl := range wls {
		// 检查类型是否正确
		if wl.Type != "allow" && wl.Type != "deny" {
			return gerror.NewErr(fmt.Sprintf("invalid type of white list: %s", wl.Type))
		}
		// 检查target是否为ip、网段。
		if !utils.IsValidNetSeg(wl.Target) && !utils.IsValidIp(wl.Target) && wl.Target != "all" {
			return gerror.NewErr(fmt.Sprintf("invalid target of white list: %s", wl.Target))
		}
	}
	return nil
}

func (n *Nginx) CheckConfigL4(l4s []NginxConfL4) *gerror.Gerr {
	for _, l4 := range l4s {
		// 判断监听端口是否合法
		if !utils.IsValidPort(l4.Listen) {
			return gerror.NewErr(fmt.Sprintf("invalid port: %d", l4.Listen))
		}
		// 校验限速语法
		re := "^[0-9]+[bBmMgG]?$"
		regular := regexp.MustCompile(re)
		if !(l4.ProxyUploadRate == "" || regular.MatchString(l4.ProxyUploadRate)) {
			return gerror.NewErr(fmt.Sprintf("invalid value of variable ProxyUploadRate: %s", l4.ProxyUploadRate))
		}
		if !(l4.ProxyDownloadRate == "" || regular.MatchString(l4.ProxyDownloadRate)) {
			return gerror.NewErr(fmt.Sprintf("invalid value of variable ProxyDownloadRate: %s", l4.ProxyDownloadRate))
		}
		// 校验超时语法
		re = "^[0-9]+[smh]?$"
		regular = regexp.MustCompile(re)
		if !(l4.ProxyConnectTimeout == "" || regular.MatchString(l4.ProxyConnectTimeout)) {
			return gerror.NewErr(fmt.Sprintf("invalid value of variable ProxyConnectTimeout: %s", l4.ProxyConnectTimeout))
		}
		if !(l4.ProxyTimeout == "" || regular.MatchString(l4.ProxyTimeout)) {
			return gerror.NewErr(fmt.Sprintf("invalid value of variable ProxyTimeout: %s", l4.ProxyTimeout))
		}
		// 参数校验：upstream
		for _, host := range l4.Upstream.Hosts {
			// 判断上游端口是否合法
			if !utils.IsValidPort(host.Port) {
				return gerror.NewErr(fmt.Sprintf("invalid port: %d", host.Port))
			}
			// 校验参数：host
			if !utils.IsValidIp(host.Ip) {
				return gerror.NewErr(fmt.Sprintf("invalid ip: %s", host.Ip))
			}
		}
		// 校验参数：白名单
		gerr := n.IsValidWhiteList(l4.WhiteList)
		if gerr != nil {
			return gerr
		}
		// 参数校验：日志
		reLogBuffer := "^[0-9]+[km]?$"
		regular = regexp.MustCompile(reLogBuffer)
		if l4.Log.Buffer != "" && !regular.MatchString(l4.Log.Buffer) {
			return gerror.NewErr(fmt.Sprintf("invalid buffer: %s,Expected: %s", l4.Log.Buffer, reLogBuffer))
		}
		reLogFlush := "^[0-9]+[sm]?$"
		regular = regexp.MustCompile(reLogFlush)
		if l4.Log.Flush != "" && !regular.MatchString(l4.Log.Flush) {
			return gerror.NewErr(fmt.Sprintf("invalid flush: %s,Expected: %s", l4.Log.Flush, reLogFlush))
		}
	}
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
	// 生成注释
	var comments []string
	for _, comment := range l4.Comments {
		comments = append(comments, "# "+comment)
	}
	commentStr := strings.Join(comments, "\n")
	// 开始生成server块
	var confServer []string
	// 生成upstream
	confUpstream := n.genConfigL4Upstream(l4.Upstream, l4.Listen)
	// 生成server语句
	confServer = append(confServer, fmt.Sprintf("listen %d;", l4.Listen))
	// 生成access_log语句
	confLog := n.genConfigL4Log(l4)
	if confLog != "" {
		confServer = append(confServer, confLog)
	}
	// 生成proxy_pass语句
	confServer = append(confServer, fmt.Sprintf("proxy_pass %d;", l4.Listen))
	// 限速
	if l4.ProxyUploadRate != "" {
		confServer = append(confServer, fmt.Sprintf("proxy_upload_rate %s;", l4.ProxyUploadRate))
	}
	if l4.ProxyDownloadRate != "" {
		confServer = append(confServer, fmt.Sprintf("proxy_download_rate %s;", l4.ProxyDownloadRate))
	}
	// 超时
	if l4.ProxyConnectTimeout != "" {
		confServer = append(confServer, fmt.Sprintf("proxy_connect_timeout %s;", l4.ProxyConnectTimeout))
	}
	if l4.ProxyTimeout != "" {
		confServer = append(confServer, fmt.Sprintf("proxy_timeout %s;", l4.ProxyTimeout))
	}
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
	return fmt.Sprintf("%s\n%s\nserver {\n\t%s\n}", commentStr, confUpstream, strings.Join(confServer, "\n\t")), nil
}
func (n *Nginx) genConfigL4Log(l4 NginxConfL4) string {
	if l4.Log.Mod == "local" {
		// 日志路径
		if l4.Log.Path == "" {
			l4.Log.Path = filepath.Join(conf.GConf.L4.Log.Path, fmt.Sprintf("%d.stream.log", l4.Listen))
		}
		ret := fmt.Sprintf("access_log %s", l4.Log.Path)
		// 日志格式
		if l4.Log.FormatName == "" {
			l4.Log.FormatName = conf.GConf.L4.Log.FormatNameDefault
		}
		if l4.Log.FormatName != "" {
			ret = ret + " " + l4.Log.FormatName
		}
		// buffer
		if l4.Log.Buffer != "" {
			ret = ret + " " + fmt.Sprintf("buffer=%s", l4.Log.Buffer)
		}
		// flush
		if l4.Log.Flush != "" {
			ret = ret + " " + fmt.Sprintf("flush=%s", l4.Log.Flush)
		}
		return ret + ";"
	} else if l4.Log.Mod == "off" {
		return "access_log off;"
	} else {
		return ""
	}
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
	var gerr *gerror.Gerr
	if l4.Disable {
		gerr = n.delFile(l4)
	} else {
		gerr = n.writeFile(l4)
	}
	if gerr != nil {
		return gerr
	}
	_, gerr = n.testAndReload()
	return gerr
}
func (n *Nginx) delFile(l4 NginxConfL4) *gerror.Gerr {
	targetFile := n.GenFilePathL4(l4)
	err := utils.DelFile(targetFile)
	if err != nil {
		return gerror.NewErr(err.Error())
	}
	return nil
}
func (n *Nginx) delFileAndReload(l4 NginxConfL4) *gerror.Gerr {
	gerr := n.delFile(l4)
	if gerr != nil {
		return gerr
	}
	_, gerr = n.testAndReload()
	return gerr
}
func (n *Nginx) Test() (string, *gerror.Gerr) {
	// fmt.Println("发起nginx测试")
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), gerror.NewErr(string(output) + "." + err.Error())
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
	output, gerr := n.Test()
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
		return gerror.NewErr(err.Error())
	}
	// 写入配置
	var gerrs []string
	for _, conf := range AllConfig {
		gerr := n.writeFileAndReload(conf)
		if gerr != nil {
			gerrs = append(gerrs, gerr.Error())
		}
	}
	if len(gerrs) > 0 {
		return gerror.NewErr(strings.Join(gerrs, ","))
	}
	return nil
}

// Get 返回值中，监听端口为0，表示配置不存在
func (n *Nginx) Get(listenPort int) (NginxConfL4, *gerror.Gerr) {
	kv, gerr := GEtcd.GetL4(listenPort)
	if gerr != nil {
		return NginxConfL4{}, gerr
	}
	if kv == nil {
		return NginxConfL4{}, nil
	}
	var l4 NginxConfL4
	err := json.Unmarshal(kv.Value, &l4)
	if err != nil {
		logger.GLogger.WithFields(logrus.Fields{"key": string(kv.Key), "value": string(kv.Value)}).Error("解析json失败")
	}
	return l4, nil
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
