package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func RaidLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.RaidLoginResponse)

	rsp.SeasonType = proto.RaidSeasonType_Close
}
