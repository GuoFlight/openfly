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
		Listen:              8081,
		Upstream:            upstream,
		IncludeFiles:        []string{"/a/b.conf", "/c/d.conf"},
		Comments:            []string{"我是注释", "I am the comment."},
		ProxyDownloadRate:   "30M",
		ProxyUploadRate:     "20M",
		ProxyConnectTimeout: "1s",
		ProxyTimeout:        "2m",
		Log: NginxLog{
			Mod:        "local",
			FormatName: "main",
			Path:       "",
			Buffer:     "5k",
			Flush:      "5s",
		},
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
	conf.PathData = "/tmp"
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
		Comments:     []string{"我是注释", "I am the comment."},
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

func TestNginx_Get(t *testing.T) {
	conf.ParseConfig("../config.toml")
	logger.InitLog()
	got, gerr := GNginx.Get(30001)
	if gerr != nil {
		t.Fatal(gerr)
	}
	fmt.Println(got)
}

func TestNginx_CheckConfigL4(t *testing.T) {
	gerr := GNginx.CheckConfigL4([]NginxConfL4{
		{
			Listen:          80,
			ProxyUploadRate: "12M",
			ProxyTimeout:    "11h",
		},
	})
	if gerr != nil {
		fmt.Println(gerr)
	}
}

func TestNginx_genConfigL4Log(t *testing.T) {
	l4 := NginxConfL4{Listen: 8081}
	// 关闭日志
	l4.Log = NginxLog{
		Mod:        "off",
		FormatName: "",
		Path:       "",
		Buffer:     "",
		Flush:      "",
	}
	conf := GNginx.genConfigL4Log(l4)
	fmt.Println(conf)
	// 输出空
	l4.Log = NginxLog{
		Mod: "global",
	}
	conf = GNginx.genConfigL4Log(l4)
	fmt.Println(conf)
	// 输出空
	l4.Log = NginxLog{
		Mod: "hosydfoshdfw0efwef",
	}
	conf = GNginx.genConfigL4Log(l4)
	fmt.Println(conf)
	// 自动生成日志路径
	l4.Log = NginxLog{
		Mod:  "local",
		Path: "",
	}
	conf = GNginx.genConfigL4Log(l4)
	fmt.Println(conf)
	// 指定日志路径和log_format
	l4.Log = NginxLog{
		Mod:        "local",
		Path:       "/var/log/nginx/test.log",
		FormatName: "test_log",
	}
	conf = GNginx.genConfigL4Log(l4)
	fmt.Println(conf)
	// 指定buffer和flush
	l4.Log = NginxLog{
		Mod:        "local",
		Path:       "/var/log/nginx/test.log",
		FormatName: "test_log",
		Buffer:     "5k",
		Flush:      "5s",
	}
	conf = GNginx.genConfigL4Log(l4)
	fmt.Println(conf)
}
