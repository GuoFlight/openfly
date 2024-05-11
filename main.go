package main

import (
	"fmt"
	"openfly/api"
	"openfly/common"
	"openfly/conf"
	"openfly/flag"
	"openfly/logger"
)

func main() {
	// 解析命令行
	flag.InitFlag()
	if *flag.Version {
		fmt.Println(conf.Version)
		return
	}

	// 解析配置文件
	conf.ParseConfig(*flag.PathConfFile)

	// 初始化日志
	logger.InitLog()

	// 更新Nginx配置
	gerr := common.GNginx.Reset()
	if gerr != nil {
		logger.GLogger.Fatal("初始化Nginx配置失败:", gerr)
	}

	// 监控etcd
	go common.GEtcd.StartWatch()

	// 启动http服务
	go api.StartHttpServer()

	// 阻塞主进程
	select {}
}
