package gateway

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/cmd"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pack"
	"github.com/gucooing/BaPs/pkg/logger"
)

type handlerFunc func(s *enter.Session, request, response mx.Message)

func (g *Gateway) newFuncRouteMap() {
	g.funcRouteMap = map[int32]handlerFunc{
		proto.Protocol_Account_CheckYostar:        g.AccountCheckYostar,           // 验证EnterTicket
		proto.Protocol_Account_Auth:               pack.AccountAuth,               // 账号验证
		proto.Protocol_Account_Nickname:           pack.AccountNickname,           // 设置昵称
		proto.Protocol_ProofToken_RequestQuestion: pack.ProofTokenRequestQuestion, // 验证登录token
		proto.Protocol_NetworkTime_Sync:           pack.NetworkTimeSync,           // 同步时间
		proto.Protocol_Academy_GetInfo:            pack.AcademyGetInfo,            // 获取学院信息
		proto.Protocol_Account_LoginSync:          pack.AccountLoginSync,          // 同步账号信息
	}
}

func (g *Gateway) registerMessage(c *gin.Context, request mx.Message, base *mx.BasePacket) {
	// panic捕获
	defer func() {
		if err := recover(); err != nil {
			errBestHTTP(c)
			logger.Error("@LogTag(player_panic)@ cmdId:%s json:%s\nerr:%s\nstack:%s", cmd.Get().GetCmdNameByCmdId(request.GetProtocol()),
				request.String(), err, logger.Stack())
			return
		}
	}()
	handler, ok := g.funcRouteMap[request.GetProtocol()]
	if !ok {
		errBestHTTP(c)
		logger.Error("@LogTag(player_no_route)@C --> S no route for msg, cmdId: %s", cmd.Get().GetCmdNameByCmdId(request.GetProtocol()))
		return
	}
	response := cmd.Get().GetResponsePacketByCmdId(request.GetProtocol())
	if response == nil {
		errBestHTTP(c)
		logger.Debug("response unknown cmd id: %v\n", cmd.Get().GetCmdNameByCmdId(request.GetProtocol()))
		return
	}
	sessionKey := base.GetSessionKey()
	var s *enter.Session
	if sessionKey == nil &&
		request.GetProtocol() != proto.Protocol_Account_CheckYostar {
		errBestHTTP(c)
		logger.Debug("get request sessionKey nil")
		return
	} else if request.GetProtocol() != proto.Protocol_Account_CheckYostar {
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
