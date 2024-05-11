package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GuoFlight/gerror"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"openfly/conf"
	"openfly/logger"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type Etcd struct {
}

var GEtcd Etcd

func (e Etcd) Connect() (*clientv3.Client, error) {
	// 建立连接
	config := clientv3.Config{
		Endpoints: []string{conf.GConf.Etcd.Server},
	}
	client, err := clientv3.New(config)
	if err != nil {
		return client, err
	}
	// 检测是否连接成功
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(conf.GConf.Etcd.Timeout))
	defer cancel()
	_, err = client.Status(timeoutCtx, config.Endpoints[0])
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Get 获取某个Key
// 不存在则返回nil
func (e Etcd) Get(key string) (*mvccpb.KeyValue, *gerror.Gerr) {
	// 建立连接
	connect, err := e.Connect()
	if err != nil {
		return nil, gerror.NewErr(err.Error())
	}
	kv := clientv3.NewKV(connect)
	// 查询
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.GConf.Etcd.Timeout)*time.Second)
	defer cancel()
	res, err := kv.Get(ctx, key)
	if err != nil {
		return nil, gerror.NewErr(err.Error())
	}
	if len(res.Kvs) == 0 {
		return nil, nil
	}
	return res.Kvs[0], nil
}
func (e Etcd) Write(k, v string) *gerror.Gerr {
	// 建立连接
	connect, err := e.Connect()
	if err != nil {
		return gerror.NewErr(err.Error())
	}
	kv := clientv3.NewKV(connect)
	// 写入
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.GConf.Etcd.Timeout)*time.Second)
	defer cancel()
	_, err = kv.Put(ctx, k, v)
	if err != nil {
		return gerror.NewErr(err.Error())
	}
	return nil
}

// Delete 删除Key
func (e Etcd) Delete(k string) *gerror.Gerr {
	// 建立连接
	connect, err := e.Connect()
	if err != nil {
		return gerror.NewErr(err.Error())
	}
	kv := clientv3.NewKV(connect)
	// 检查是否存在此Key
	kvL4, gerr := GEtcd.Get(k)
	if gerr != nil {
		return gerr
	}
	if kvL4 == nil {
		return gerror.NewErr(fmt.Sprintf("This key does not exist：%s", k))
	}
	// 删除
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.GConf.Etcd.Timeout)*time.Second)
	defer cancel()
	_, err = kv.Delete(ctx, k)
	if err != nil {
		return logger.PrintErr(gerror.NewErr(err.Error()), nil)
	}
	return nil
}
func (e Etcd) StartWatch() {
	connect, err := e.Connect()
	if err != nil {
		logger.GLogger.Error(err)
		os.Exit(1)
	}

	watcher := connect.Watch(context.Background(), conf.GConf.Etcd.Prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for watchResponse := range watcher {
		for _, event := range watchResponse.Events {
			var l4 NginxConfL4
			switch event.Type {
			case mvccpb.PUT:
				if event.PrevKv == nil {
					logger.GLogger.Infof("监控到Key变化!Type=ADD Key=%s NewValue=%s\n", event.Kv.Key, event.Kv.Value)
				} else {
					logger.GLogger.Infof("监控到Key变化!Type=UPDATE Key=%s OldValue=%s NewValue=%s\n", event.Kv.Key, event.PrevKv.Value, event.Kv.Value)
				}
				err := json.Unmarshal(event.Kv.Value, &l4)
				if err != nil {
					logger.GLogger.Error("无法解析etcd中的json:", string(event.Kv.Value))
					continue
				}
				gerr := GNginx.WriteFileAndReload(l4)
				if gerr != nil {
					logger.GLogger.Error("写入文件错误:", gerr)
					continue
				}
			case mvccpb.DELETE:
				logger.GLogger.Infof("监控到Key变化!Type=DELETE Key=%s OldValue=%s", event.Kv.Key, event.PrevKv.Value)
				err := json.Unmarshal(event.PrevKv.Value, &l4)
				if err != nil {
					logger.GLogger.Error("无法解析etcd中的json:", string(event.PrevKv.Value))
					continue
				}
				gerr := GNginx.DelFileAndReload(l4)
				if gerr != nil {
					logger.GLogger.Error("无法删除文件:", gerr)
					continue
				}
			default:
				logger.GLogger.Warn("Unknown etcd event type!")
			}
		}
	}
}

func (e Etcd) GenKeyL4(listen int) string {
	return filepath.Join(conf.GConf.Etcd.Prefix, conf.EtcdSubPathL4, strconv.Itoa(listen))
}
func (e Etcd) GetL4(listen int) (*mvccpb.KeyValue, *gerror.Gerr) {
	return e.Get(e.GenKeyL4(listen))
}
func (e Etcd) WriteL4(l4 NginxConfL4) *gerror.Gerr {
	l4json, err := json.Marshal(l4)
	if err != nil {
		return gerror.NewErr(err.Error())
	}
	gerr := e.Write(e.GenKeyL4(l4.Listen), string(l4json))
	if gerr != nil {
		return gerr
	}
	return nil
}

// GetAllL4 获取所有L4配置
func (e Etcd) GetAllL4() ([]*mvccpb.KeyValue, *gerror.Gerr) {
	// 建立连接
	connect, err := e.Connect()
	if err != nil {
		return nil, gerror.NewErr(err.Error())
	}
	kv := clientv3.NewKV(connect)
	// 查询
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.GConf.Etcd.Timeout)*time.Second)
	defer cancel()
	res, err := kv.Get(ctx, path.Join(conf.GConf.Etcd.Prefix, conf.EtcdSubPathL4), clientv3.WithPrefix())
	if err != nil {
		return nil, gerror.NewErr(err.Error())
	}
	return res.Kvs, nil
}
func (e Etcd) DeleteL4(listen int) *gerror.Gerr {
	return e.Delete(e.GenKeyL4(listen))
}
func (e Etcd) AddL4(l4 NginxConfL4) *gerror.Gerr {
	kv, gerr := e.GetL4(l4.Listen)
	if gerr != nil {
		return gerr
	}
	if kv != nil {
		return gerror.NewErr(fmt.Sprintf("This configuration already exists：%d", l4.Listen))
	}
	// 添加配置
	return e.WriteL4(l4)
}
