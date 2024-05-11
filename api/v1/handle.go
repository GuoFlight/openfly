package v1

import (
	"github.com/gin-gonic/gin"
)

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Health(c *gin.Context) {
	c.String(200, "alive")
}
