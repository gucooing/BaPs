package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func MailCheck(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MailCheckResponse)

	rsp.Count = 0 // 未领取数量
}
