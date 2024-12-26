package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/mx/protocol"
	"github.com/gucooing/BaPs/pkg/logger"
)

func NewGateWay(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/gateway", gateway)
	}
}

func gateway(c *gin.Context) {
	if c.GetHeader("user-agent") != "BestHTTP/2 v2.4.0" ||
		c.GetHeader("accept-encoding") != "gzip" {
		return
	}
	file, err := c.FormFile("mx")
	if err != nil {
		logger.Debug("", err.Error())
		return
	}
	fileContent, err := file.Open()
	if err != nil {
		logger.Debug("", err.Error())
		return
	}
	bin := make([]byte, file.Size)
	_, err = fileContent.Read(bin)
	if err != nil {
		return
	}
	s, err := protocol.Decode(bin)
	logger.Debug("c--->s:%s", s)
	// YostarLoginToken 只能用一次！！！！！！！！！
}
