package v1

import (
	"github.com/gin-gonic/gin"
	"openfly/common"
)

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Health(c *gin.Context) {
	_, gerr := common.GNginx.Test()
	if gerr != nil {
		c.JSON(200, Res{
			Code: 1,
			Msg:  gerr.Error(),
		})
		return
	}
	c.JSON(200, Res{
		Code: 0,
		Msg:  "healthy",
	})
}
