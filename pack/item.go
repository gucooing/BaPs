package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func ItemList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ItemListResponse)

	rsp.ExpiryItemDBs = make([]*proto.ItemDB, 0)
	rsp.ItemDBs = make([]*proto.ItemDB, 0)

	rsp.ItemDBs = append(rsp.ItemDBs, &proto.ItemDB{
		Type:       proto.ParcelType_Item,
		ServerId:   game.GetServerId(),
		UniqueId:   2,
		StackCount: 5,
	})
}

func EquipmentList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.EquipmentItemListResponse)

	rsp.EquipmentDBs = make([]*proto.EquipmentDB, 0)
}
