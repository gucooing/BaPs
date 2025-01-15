package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/logger"
)

func ShopGachaRecruitList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ShopGachaRecruitListResponse)

	rsp.ShopRecruits = make([]*proto.ShopRecruitDB, 0)                         // 卡池数据
	rsp.ShopFreeRecruitHistoryDBs = make([]*proto.ShopFreeRecruitHistoryDB, 0) // 免费抽卡历史数据
}

func ShopBeforehandGachaGet(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ShopBeforehandGachaGetResponse)

	bin := game.GetBeforehandInfo(s)
	if bin == nil {
		return
	}

	rsp.AlreadyPicked = bin.AlreadyPicked
}

func ShopBeforehandGachaRun(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopBeforehandGachaRunRequest)
	rsp := response.(*proto.ShopBeforehandGachaRunResponse)

	bin := game.GetBeforehandInfo(s)
	if bin == nil {
		return
	}
	bin.GoodsId = req.GoodsId
	bin.ShopUniqueId = req.ShopUniqueId

	if bin.LastIndex < 10 {
		bin.LastResults = game.GachaRun(10, true, true)
	} else {
		bin.LastIndex-- // 避免傻子客户端溢出
	}
	rsp.SelectGachaSnapshot = game.GetBeforehandGachaSnapshotDB(s)
	bin.LastIndex++
}

func ShopBeforehandGachaSave(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopBeforehandGachaSaveRequest)
	rsp := response.(*proto.ShopBeforehandGachaSaveResponse)

	bin := game.GetBeforehandInfo(s)
	if bin == nil {
		return
	}
	if bin.LastIndex-1 == req.TargetIndex {
		bin.SavedIndex = req.TargetIndex
		bin.SavedResults = bin.LastResults
	}
	rsp.SelectGachaSnapshot = game.GetBeforehandGachaSnapshotDB(s)
}

func ShopBeforehandGachaPick(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopBeforehandGachaPickRequest)
	rsp := response.(*proto.ShopBeforehandGachaPickResponse)

	rsp.GachaResults = make([]*proto.GachaResult, 0)
	rsp.AcquiredItems = make([]*proto.ItemDB, 0)
	bin := game.GetBeforehandInfo(s)
	if bin == nil {
		return
	}
	var results []int64
	if bin.LastIndex-1 == req.TargetIndex {
		results = bin.LastResults
	} else if bin.SavedIndex == req.TargetIndex {
		results = bin.SavedResults
	} else {
		logger.Debug("[UID:%v]新手免费十连,错误的抽卡记录", s.AccountServerId)
		return
	}

	bin.AlreadyPicked = true
	list, addItemList := game.SaveGachaResults(s, results)
	rsp.GachaResults = list
	for id, num := range addItemList {
		itemInfo := game.GetItemInfo(s, id)
		if itemInfo == nil {
			continue
		}
		rsp.AcquiredItems = append(rsp.AcquiredItems, &proto.ItemDB{
			Type:       proto.ParcelType_Item,
			ServerId:   itemInfo.ServerId,
			UniqueId:   id,
			StackCount: num,
		})
	}
}
