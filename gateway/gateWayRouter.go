package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/pkg/logger"
)

type handlerFunc func(s *enter.Session, packet mx.Message) mx.Message

var handlerFuncRouteMap = map[int32]handlerFunc{}

func (g *Gateway) registerMessage(c *gin.Context, packet mx.Message) {
	// panic捕获
	defer func() {
		if err := recover(); err != nil {
			logger.Error("@LogTag(player_panic)@ cmdId:%s b64:%s json:%s\nerr:%s\nstack:%s", packet.GetProtocolValue(),
				packet.String(), err, logger.Stack())
			return
		}
	}()
	handler, ok := handlerFuncRouteMap[packet.GetProtocolKey()]
	if !ok {
		logger.Error("@LogTag(player_no_route)@C --> S no route for msg, cmdId: %s", packet.GetProtocolValue())
		return
	}
	rsp := handler(nil, packet)
	g.send(c, rsp)
	return
}
