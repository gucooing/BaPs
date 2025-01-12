package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func ShopBeforehandGachaGet(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ShopBeforehandGachaGetResponse)

	rsp.ServerNotification = proto.ServerNotificationFlag_HasUnreadMail
}

func ShopGachaRecruitList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ShopGachaRecruitListResponse)

	rsp.ShopRecruits = make([]*proto.ShopRecruitDB, 0)
	rsp.ShopFreeRecruitHistoryDBs = make([]*proto.ShopFreeRecruitHistoryDB, 0)
}
