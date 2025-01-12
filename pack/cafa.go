package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func CafeGetInfo(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CafeGetInfoResponse)

	rsp.CafeDBs = game.GetPbCafeDBs(s)
	rsp.FurnitureDBs = make([]*proto.FurnitureDB, 0) // 已获得家具数据
}
