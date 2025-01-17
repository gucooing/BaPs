package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func MailCheck(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MailCheckResponse)

	rsp.Count = 0 // 未领取数量
}
