package common

import (
	"fmt"
	"openfly/conf"
	"path"
	"testing"
)

func TestEtcd_Write(t *testing.T) {
	conf.ParseConfig("../config.toml")
	gerr := GEtcd.Write(conf.GConf.Etcd.Prefix+"/test1", "test2")
	if gerr != nil {
		t.Error(gerr)
		return
	}
}

func TestEtcd_GetL4(t *testing.T) {
	conf.ParseConfig("../config.toml")
	kvs, gerr := GEtcd.GetAllL4()
	if gerr != nil {
		t.Error(gerr)
		return
	}
	for _, kv := range kvs {
		fmt.Println(path.Base(string(kv.Key)), string(kv.Value))
	}
}

func TestEtcd_Get(t *testing.T) {
	conf.ParseConfig("../config.toml")
	kv, gerr := GEtcd.Get("30001")
	if gerr != nil {
		t.Error(gerr)
		return
	}
	fmt.Println(string(kv.Key), string(kv.Value))
}
