package common

import (
	"fmt"
	"openfly/conf"
	"openfly/logger"
	"testing"
)

func TestNginx_genConfigL4Upstream(t *testing.T) {
	conf.ParseConfig("../config.toml")
	upstream := Upstream{
		Hosts: []UpstreamHost{
			{
				Ip:   "1.1.1.1",
				Port: 80,
			},
			{
				Ip:       "2.2.2.2",
				Port:     81,
				IsBackup: true,
			},
		},
		IsHash:    true,
		HashField: "test_field",
		Interval:  10,
		Rise:      2,
		Fall:      2,
		Timeout:   1000,
	}
	confUpstream := GNginx.genConfigL4Upstream(upstream, 8888)
	fmt.Println(confUpstream)
}
func TestNginx_GenConfigL4(t *testing.T) {
	conf.ParseConfig("../config.toml")
	upstream := Upstream{
		Hosts: []UpstreamHost{
			{
				Ip:   "1.1.1.1",
				Port: 80,
			},
			{
				Ip:       "2.2.2.2",
				Port:     81,
				IsBackup: true,
			},
		},
		IsHash:    true,
		HashField: "test_field",
		Interval:  10,
		Rise:      2,
		Fall:      2,
		Timeout:   1000,
	}
	confL4 := NginxConfL4{
		Listen:       8081,
		Upstream:     upstream,
		IncludeFiles: []string{"/a/b.conf", "/c/d.conf"},
		WhiteList: []WhiteListItem{
			{
				conf.Allow,
				"1.1.1.1",
			},
			{
				conf.Allow,
				"2.2.2.2",
			},
			{
				conf.Deny,
				"all",
			},
		},
	}
	fmt.Println(GNginx.GenConfigL4(confL4))
}
func TestNginx_WriteFile(t *testing.T) {
	conf.ParseConfig("../config.toml")
	conf.PathData = "../test/l4"
	upstream := Upstream{
		Hosts: []UpstreamHost{
			{
				Ip:   "1.1.1.1",
				Port: 80,
			},
			{
				Ip:       "2.2.2.2",
				Port:     81,
				IsBackup: true,
			},
		},
		IsHash:    true,
		HashField: "test_field",
		Interval:  10,
		Rise:      2,
		Fall:      2,
		Timeout:   1000,
	}
	confL4 := NginxConfL4{
		Listen:       8081,
		Upstream:     upstream,
		IncludeFiles: []string{"/a/b.conf", "/c/d.conf"},
		WhiteList: []WhiteListItem{
			{
				conf.Allow,
				"1.1.1.1",
			},
			{
				conf.Allow,
				"2.2.2.2",
			},
			{
				conf.Deny,
				"all",
			},
		},
	}
	gerr := GNginx.WriteFileAndReload(confL4)
	if gerr != nil {
		t.Error(gerr)
		return
	} else {
		fmt.Println("文件生成在:", conf.PathData)
	}
}
func TestNginx_GetAll(t *testing.T) {
	conf.ParseConfig("../config.toml")
	logger.InitLog()
	l4s, gerr := GNginx.GetAll()
	if gerr != nil {
		t.Error(gerr)
		return
	}
	for _, v := range l4s {
		fmt.Println(v.Listen)
	}
}
