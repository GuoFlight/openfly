package v1

import (
	"fmt"
	"github.com/GuoFlight/gerror"
	"github.com/gin-gonic/gin"
	"openfly/common"
	"openfly/logger"
	"strconv"
)

func Set(c *gin.Context) {
	var req common.NginxConfL4
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, Res{
			Code: 400,
			Msg:  "invalid request." + err.Error(),
		})
		return
	}
	// 参数校验
	gerr := common.GNginx.CheckConfigL4([]common.NginxConfL4{req})
	if gerr != nil {
		c.JSON(400, Res{
			Code: 400,
			Msg:  gerr.Error(),
		})
		return
	}
	// 更新配置
	gerr = common.GEtcd.WriteL4(req)
	if gerr != nil {
		logger.PrintErr(gerr, nil)
		c.JSON(500, Res{
			Code: 500,
			Msg:  gerr.Error(),
		})
		return
	}
	c.JSON(0, Res{
		Code: 0,
		Msg:  "success",
	})
}
func Add(c *gin.Context) {
	var req common.NginxConfL4
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, Res{
			Code: 400,
			Msg:  "invalid request." + err.Error(),
		})
		return
	}
	// 参数校验
	gerr := common.GNginx.CheckConfigL4([]common.NginxConfL4{req})
	if gerr != nil {
		c.JSON(400, Res{
			Code: 400,
			Msg:  gerr.Error(),
		})
		return
	}
	// 添加配置
	gerr = common.GEtcd.AddL4(req)
	if gerr != nil {
		logger.PrintErr(gerr, nil)
		c.JSON(400, Res{
			Code: 400,
			Msg:  gerr.Error(),
		})
		return
	}
	c.JSON(0, Res{
		Code: 0,
		Msg:  "success",
	})
}
func Get(c *gin.Context) {
	listen := c.DefaultQuery("listen", "")
	if listen == "" {
		GetAll(c)
		return
	}
	listenPort, err := strconv.Atoi(listen)
	if err != nil {
		logger.PrintErr(gerror.NewErr(err.Error()), nil)
		c.JSON(400, Res{
			Code: 400,
			Msg:  fmt.Sprintf("invalid port：%s", listen),
		})
		return
	}
	l4, gerr := common.GNginx.Get(listenPort)
	if gerr != nil {
		c.JSON(500, Res{
			Code: 500,
			Msg:  gerr.Error(),
		})
	}
	if l4.Listen == 0 {
		c.JSON(404, Res{
			Code: 404,
			Msg:  fmt.Sprintf("Listening port does not exist: %d", listenPort),
		})
		return
	}
	c.JSON(0, Res{
		Code: 0,
		Data: l4,
	})
}
func GetAll(c *gin.Context) {
	l4s, gerr := common.GNginx.GetAll()
	if gerr != nil {
		logger.PrintErr(gerr, nil)
		c.JSON(500, Res{
			Code: 500,
			Msg:  gerr.Error(),
		})
		return
	}
	c.JSON(0, Res{
		Code: 0,
		Data: l4s,
	})
}
func Delete(c *gin.Context) {
	// 获取参数
	listenStr := c.Query("listen")
	if listenStr == "" {
		c.JSON(400, Res{
			Code: 400,
			Msg:  "Missing parameter listen",
		})
		return
	}
	listen, err := strconv.Atoi(listenStr)
	if err != nil {
		c.JSON(400, Res{
			Code: 400,
			Msg:  "Parameter listen is not a number",
		})
		return
	}
	// 删除配置
	gerr := common.GEtcd.DeleteL4(listen)
	if gerr != nil {
		c.JSON(200, Res{
			Code: 1,
			Msg:  gerr.Error(),
		})
		return
	}
	c.JSON(200, Res{
		Code: 0,
		Msg:  "success",
	})
}
