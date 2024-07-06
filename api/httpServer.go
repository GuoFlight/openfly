package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "openfly/api/v1"
	"openfly/conf"
	"openfly/logger"
)

func StartHttpServer() {
	router := gin.New()
	router.Use(gin.Recovery())

	// 子路由v1
	subRouterV1 := router.Group("/v1")
	subRouterV1.GET("/health", v1.Health)
	subRouterV1.POST("/login", v1.Login)

	// 子路由/v1/admin
	subRouterV1Admin := subRouterV1.Group("/admin")
	subRouterV1Admin.Use(v1.Jwt())
	// 子路由/v1/admin/nginx
	subRouterV1AdminNginx := subRouterV1Admin.Group("/nginx")
	subRouterV1AdminNginx.POST("/set", v1.Set)
	subRouterV1AdminNginx.POST("/add", v1.Add)
	subRouterV1AdminNginx.GET("/get", v1.Get)
	subRouterV1AdminNginx.GET("/getAll", v1.GetAll)
	subRouterV1AdminNginx.DELETE("/delete", v1.Delete)
	subRouterV1AdminNginx.POST("/switch", v1.Switch)

	// 启动http服务
	logger.GLogger.Info("开始启动http服务")
	err := router.Run(fmt.Sprintf(":%d", conf.GConf.Http.Port))
	if err != nil {
		logger.GLogger.Fatal(err)
	}
}
