package gateway

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol"
	"github.com/gucooing/BaPs/protocol/proto"
)

type Gateway struct {
	router *gin.Engine
}

func NewGateWay(router *gin.Engine) *Gateway {
	g := &Gateway{
		router: router,
	}
	enter.MaxCachePlayerTime = alg.MaxInt(config.GetGateWay().MaxCachePlayerTime, 30)
	enter.MaxPlayerNum = config.GetGateWay().MaxPlayerNum
	g.initRouter()

	return g
}

func (g *Gateway) initRouter() {
	g.router.POST("/getEnterTicket/gateway", g.getEnterTicket)
	api := g.router.Group("/api")
	{
		api.POST("/gateway", g.gateWay)
	}
}

func (g *Gateway) send(c *gin.Context, n proto.Message) {
	rsp, err := protocol.MarshalResponse(n)
	if err != nil {
		logger.Debug("marshal err:", err)
		return
	}
	c.JSON(200, rsp)
}

func (g *Gateway) gateWay(c *gin.Context) {
	if !alg.CheckGateWay(c) {
		return
	}
	bin, err := mx.GetFormMx(c)
	if err != nil {
		logger.Warn("get form mx error:", err)
		return
	}
	packet, base, err := protocol.UnmarshalRequest(bin)
	if err != nil {
		errBestHTTP(c, 15022)
		logger.Debug("unmarshal c--->s err:%s,json:%s", err, string(bin))
		return
	}
	g.registerMessage(c, packet, base)
}

func errBestHTTP(c *gin.Context, errorCode int32) {
	msg := &protocol.NetworkProtocolResponse{
		Packet:   fmt.Sprintf("{\"Protocol\":-1,\"ErrorCode\":%d}", errorCode),
		Protocol: "Error",
	}
	c.JSON(200, msg)
}

func errTokenBestHTTP(c *gin.Context) {
	c.JSON(200, gin.H{
		"protocol": "Error",
		"packet":   "{\"Protocol\":-1,\"ErrorCode\":500}",
	})
}
