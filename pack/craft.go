package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func CraftInfoList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CraftInfoListResponse)

	rsp.CraftInfos = make([]*proto.CraftInfoDB, 0)
	rsp.ShiftingCraftInfos = make([]*proto.ShiftingCraftInfoDB, 0)
}
