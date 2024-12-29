package gateway

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/cmd"
	"github.com/gucooing/BaPs/pkg/logger"
)

type handlerFunc func(s *enter.Session, request, response mx.Message)

func (g *Gateway) newFuncRouteMap() {
	g.funcRouteMap = map[int32]handlerFunc{
		mx.Protocol_Account_CheckYostar: g.AccountCheckYostar, // 验证EnterTicket
		mx.Protocol_Account_Auth:        game.AccountAuth,     // 账号验证
	}
}

func (g *Gateway) registerMessage(c *gin.Context, request mx.Message, base *mx.BasePacket) {
	// panic捕获
	defer func() {
		if err := recover(); err != nil {
			errBestHTTP(c)
			logger.Error("@LogTag(player_panic)@ cmdId:%s json:%s\nerr:%s\nstack:%s", request.GetProtocolValue(),
				request.String(), err, logger.Stack())
			return
		}
	}()
	handler, ok := g.funcRouteMap[request.GetProtocolKey()]
	if !ok {
		errBestHTTP(c)
		logger.Error("@LogTag(player_no_route)@C --> S no route for msg, cmdId: %s", request.GetProtocolValue())
		return
	}
	response := cmd.Get().GetResponsePacketByCmdId(request.GetProtocolKey())
	if response == nil {
		errBestHTTP(c)
		logger.Debug("response unknown cmd id: %v\n", request.GetProtocolValue())
		return
	}
	sessionKey := base.GetSessionKey()
	var s *enter.Session
	if sessionKey == nil &&
		request.GetProtocolKey() != mx.Protocol_Account_CheckYostar {
		errBestHTTP(c)
		logger.Debug("get request sessionKey nil")
		return
	} else if request.GetProtocolKey() != mx.Protocol_Account_CheckYostar {
		s = enter.GetSessionBySessionKey(sessionKey)
		if s == nil {
			errBestHTTP(c)
			logger.Debug("get session nil")
			return
		}
		s.EndTime = time.Now().Add(30 * time.Minute)
	}
	response.SetSessionKey(base)
	handler(s, request, response)
	logger.Debug("gateway c--->s :%s", request.String())
	logger.Debug("gateway s--->c :%s", response.String())
	g.send(c, response)
	return
}

// 62135629200 * 10000000 + curTime
