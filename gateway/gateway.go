package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/protocol"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
)

type Gateway struct {
	router       *gin.Engine
	snow         *alg.SnowflakeWorker
	funcRouteMap map[int32]handlerFunc
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
	g.newFuncRouteMap()
	g.router.POST("/getEnterTicket/gateway", g.getEnterTicket)
	api := g.router.Group("/api")
	{
		api.POST("/gateway", g.gateWay)
	}
}

func (g *Gateway) send(c *gin.Context, n mx.Message) {
	rsp, err := protocol.MarshalResponse(n)
	if err != nil {
		logger.Debug("marshal err:", err)
		return
	}
	c.JSON(200, rsp)
}

func (g *Gateway) gateWay(c *gin.Context) {
	if !alg.CheckGateWay(c) {
		errBestHTTP(c)
		return
	}
	bin, err := mx.GetFormMx(c)
	if err != nil {
		errBestHTTP(c)
		logger.Warn("get form mx error:", err)
		return
	}
	packet, base, err := protocol.UnmarshalRequest(bin)
	if err != nil {
		errBestHTTP(c)
		logger.Debug("unmarshal c--->s err:%s,json:%s", err, string(bin))
		return
	}
	g.registerMessage(c, packet, base)
}

func errBestHTTP(c *gin.Context) {
	c.JSON(200, gin.H{
		"protocol": "Error",
		"packet":   "{\"Protocol\":-1,\"ErrorCode\":500,\"ServerTimeTicks\":114514}",
	})
}
