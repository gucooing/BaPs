package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/protocol/proto"
)

func EliminateRaidLogin(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.EliminateRaidLoginResponse)

	rsp.SeasonType = proto.RaidSeasonType_Close
}
