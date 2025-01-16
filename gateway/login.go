package gateway

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	pb "google.golang.org/protobuf/proto"
)

func (g *Gateway) getEnterTicket(c *gin.Context) {
	if !alg.CheckGateWay(c) {
		errBestHTTP(c)
		return
	}
	bin, err := mx.GetFormMx(c)
	if err != nil {
		return
	}
	rsp := &proto.QueuingGetTicketResponse{}
	defer g.send(c, rsp)
	req := new(proto.QueuingGetTicketRequest)
	err = json.Unmarshal(bin, req)
	if err != nil {
		logger.Debug("request err:%s c--->s:%s", err.Error(), string(bin))
		return
	}
	yoStarUserLogin := db.GetYoStarUserLoginByYostarUid(req.YostarUID)
	if yoStarUserLogin == nil {
		return
	}
	if yoStarUserLogin.YostarLoginToken != req.YostarToken ||
		yoStarUserLogin.YostarLoginToken == "" {
		return
	}
	yoStarUserLogin.YostarLoginToken = ""
	if err = db.UpdateYoStarUserLogin(yoStarUserLogin); err != nil {
		return
	}
	enterTicket := fmt.Sprintf("%v%s", g.snow.GenId(), alg.RandStr(10))
	if !enter.AddEnterTicket(yoStarUserLogin.AccountServerId, enterTicket) {
		return
	}
	rsp.EnterTicket = enterTicket
	rsp.SetSessionKey(&mx.BasePacket{
		Protocol: req.Protocol,
	})
	logger.Debug("EnterTicket交换成功:%s", rsp.EnterTicket)
}

func (g *Gateway) AccountCheckYostar(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountCheckYostarRequest)
	rsp := response.(*proto.AccountCheckYostarResponse)
	var err error

	tickInfo := enter.GetEnterTicketInfo(req.EnterTicket)
	if tickInfo == nil {
		rsp.ResultMessag = "EnterTicket验证失败"
		logger.Debug("EnterTicket验证失败")
		return
	}
	enter.DelEnterTicket(req.EnterTicket)
	s = enter.GetSessionByAccountServerId(tickInfo.AccountServerId)
	mxToken := fmt.Sprintf("%v%s", g.snow.GenId(), alg.RandStr(30))
	if s == nil {
		yostarGame := db.GetYostarGameByAccountServerId(tickInfo.AccountServerId)
		if yostarGame == nil {
			// new Game Player
			yostarGame, err = db.AddYostarGameByYostarUid(tickInfo.AccountServerId)
			if err != nil {
				logger.Debug("账号创建失败:%s", err.Error())
				return
			}
		}
		s = &enter.Session{
			AccountServerId: tickInfo.AccountServerId,
			PlayerBin:       new(sro.PlayerBin),
		}
		if yostarGame.BinData != nil {
			pb.Unmarshal(yostarGame.BinData, s.PlayerBin)
		} else {
			s.PlayerBin = game.NewYostarGame(tickInfo.AccountServerId)
			logger.Debug("AccountServerId:%v,新玩家登录Game,创建新账号中", tickInfo.AccountServerId)
		}
	}
	// 更新一次账号缓存
	s.MxToken = mxToken
	s.EndTime = time.Now().Add(30 * time.Minute)
	if !enter.AddSession(s) {
		logger.Debug("AccountServerId:%v,重复上线账号", tickInfo.AccountServerId)
	}
	rsp.ResultState = 1
	base := &mx.BasePacket{
		SessionKey: &mx.SessionKey{
			AccountServerId: tickInfo.AccountServerId,
			MxToken:         s.MxToken,
		},
		Protocol:           response.GetProtocol(),
		AccountId:          tickInfo.AccountServerId,
		ServerNotification: int32(game.GetServerNotification(s)),
		ServerTimeTicks:    game.GetServerTime(),
	}
	response.SetSessionKey(base)
}
