package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func ClanLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ClanLoginResponse)

	rsp.AccountClanDB = nil                        // 社团简介
	rsp.AccountClanMemberDB = &proto.ClanMemberDB{ // 本人信息
		AccountId: s.AccountServerId,
	}
}

func ClanCheck(s *enter.Session, request, response mx.Message) {}
