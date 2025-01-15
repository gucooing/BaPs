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

	for _, conf := range game.GetItemList(s) {
		rsp.ItemDBs = append(rsp.ItemDBs, &proto.ItemDB{
			Type:       proto.ParcelType_Item,
			ServerId:   conf.ServerId,
			UniqueId:   conf.UniqueId,
			StackCount: conf.StackCount,
		})
	}
}

func EquipmentList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.EquipmentItemListResponse)

	rsp.EquipmentDBs = make([]*proto.EquipmentDB, 0)
}
