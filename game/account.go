package game

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func AccountAuth(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountAuthRequest)
	rsp := response.(*proto.AccountAuthResponse)
}
