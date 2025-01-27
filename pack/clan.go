package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func ClanLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ClanLoginResponse)

	rsp.AccountClanDB = nil                        // 社团简介
	rsp.AccountClanMemberDB = &proto.ClanMemberDB{ // 本人信息
		AccountId: s.AccountServerId,
	}
	rsp.ClanAssistSlotDBs = make([]*proto.ClanAssistSlotDB, 0) // 援助信息
}

func ClanCheck(s *enter.Session, request, response mx.Message) {}
