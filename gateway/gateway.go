package gateway

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol"
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

	go status(router)

	return g
}

func (g *Gateway) initRouter() {
	g.router.POST("/getEnterTicket/gateway", g.getEnterTicket) // 这个地方要加个限速器,不然会被dos
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
	var str string
	if config.GetHttpNet().Encoding {
		byt, _ := sonic.Marshal(rsp)
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		gz.Write(byt)
		gz.Close()
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		str = b.String()
	} else {
		str, _ = sonic.MarshalString(rsp)
	}

	c.String(http.StatusOK, str)
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
		"packet":   "{\"Protocol\":-1,\"ErrorCode\":402}",
	})
}
