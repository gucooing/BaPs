package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/protocol"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
)

type Gateway struct {
	router *gin.Engine
	snow   *alg.SnowflakeWorker
}

func NewGateWay(router *gin.Engine) *Gateway {
	g := &Gateway{
		router: router,
		snow:   alg.NewSnowflakeWorker(16),
	}
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

func (g *Gateway) send(c *gin.Context, n mx.Message) {
	rsp, err := protocol.Marshal(n)
	if err != nil {
		logger.Debug("marshal err:", err)
		return
	}
	c.JSON(200, rsp)
}

func (g *Gateway) gateWay(c *gin.Context) {
	if !alg.CheckGateWay(c) {
		c.JSON(404, gin.H{})
		return
	}
	bin, err := alg.GetFormMx(c)
	if err != nil {
		c.JSON(404, gin.H{})
		logger.Warn("get form mx error:", err)
		return
	}
	packet, err := protocol.Unmarshal(bin)
	if err != nil {
		c.JSON(404, gin.H{})
		logger.Debug("unmarshal c--->s err:", err)
		return
	}
	logger.Debug("gateway c--->s :%s", string(bin))
	g.registerMessage(c, packet)
}
