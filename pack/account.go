package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func AccountAuth(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountAuthRequest)
	rsp := response.(*proto.AccountAuthResponse)

	rsp.CurrentVersion = req.Version
	rsp.AccountDB = game.GetAccountDB(s)
	rsp.AttendanceBookRewards = game.GetAttendanceBookRewards(s)

	rsp.IssueAlertInfos = make([]*proto.IssueAlertInfoDB, 0)
	rsp.StaticOpenConditions = game.GetStaticOpenConditions(s)
	game.SetLastConnectTime(s)
	s.AccountState = proto.AccountState_Normal
}

func AccountNickname(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountNicknameRequest)
	rsp := response.(*proto.AccountNicknameResponse)

	game.SetAccountNickname(s, req.Nickname)
	rsp.AccountDB = game.GetAccountDB(s)
}

func ProofTokenRequestQuestion(s *enter.Session, request, response mx.Message) {}

func NetworkTimeSync(s *enter.Session, request, response mx.Message) {}

func AccountLoginSync(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountLoginSyncRequest)
	rsp := response.(*proto.AccountLoginSyncResponse)
	for _, cmdId := range req.SyncProtocols {
		switch cmdId {
		case proto.Protocol_Cafe_Get:
			rsp.CafeGetInfoResponse = &proto.CafeGetInfoResponse{}
		case proto.Protocol_Account_CurrencySync:
			
		}
	}
}
