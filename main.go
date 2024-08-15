package main

import (
	"fmt"
	"openfly/api"
	"openfly/common"
	"openfly/conf"
	"openfly/flag"
	"openfly/logger"
	"os"
	"os/signal"
	"syscall"
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
		logger.GLogger.Error("初始化Nginx配置失败:", gerr)
	}

	// 监控etcd
	go common.GEtcd.StartWatch()

	// 启动http服务
	go api.StartHttpServer()

	// 优雅退出
	sig := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for {
			s := <-sig
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				logger.GLogger.Info("app收到退出信号：", s)
				<-api.Done
				logger.GLogger.Info("app正常退出")
				done <- true
			default:
				fmt.Println("app收到即将忽略的信号:", s)
			}
		}
	}()
	<-done
}
