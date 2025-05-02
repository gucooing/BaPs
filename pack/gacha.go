package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
	"github.com/gucooing/BaPs/protocol/proto"
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

	genLastResults := func(results []*game.ParcelResult) []int64 {
		list := make([]int64, 0)
		for _, result := range results {
			if result.ParcelType == proto.ParcelType_Character {
				list = append(list, result.ParcelId)
			}
		}
		return list
	}

	if bin.LastIndex < 10 {
		// 抽卡生成
		bin.LastResults = genLastResults(game.GenGachaResults(req.GoodsId))
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

	gen := func(results []int64) []*game.ParcelResult {
		list := make([]*game.ParcelResult, 0)
		for _, result := range results {
			list = append(list, &game.ParcelResult{
				ParcelType: proto.ParcelType_Character,
				ParcelId:   result,
				Amount:     1,
			})
		}
		return list
	}

	bin.AlreadyPicked = true
	list, addItemList := game.SaveGachaResults(s, gen(results))
	rsp.GachaResults = list
	for id, _ := range addItemList {
		itemInfo := game.GetItemInfo(s, id)
		if itemInfo == nil {
			continue
		}
		rsp.AcquiredItems = append(rsp.AcquiredItems, &proto.ItemDB{
			Type: proto.ParcelType_Item,
			ConsumableItemBaseDB: &proto.ConsumableItemBaseDB{
				Key:        nil,
				CanConsume: false,
				ServerId:   itemInfo.ServerId,
				UniqueId:   itemInfo.UniqueId,
				StackCount: int64(itemInfo.StackCount),
			},
		})
	}
}

func ShopBuyGacha3(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopBuyGacha3Request)
	rsp := response.(*proto.ShopBuyGacha3Response)

	rsp.AcquiredItems = make([]*proto.ItemDB, 0)
	// 成本计算
	rsp.ConsumedItems = GenGachaCost(s, req.Cost)

	// 抽卡生成
	results := game.GenGachaResults(req.GoodsId)

	// 生成回包
	list, addItemList := game.SaveGachaResults(s, results)
	rsp.GachaResults = list
	for id, _ := range addItemList {
		rsp.AcquiredItems = append(rsp.AcquiredItems, game.GetItemDB(s, id))
	}

	info := game.GetCurrencyInfo(s, proto.CurrencyTypes_GemBonus)
	if info != nil {
		rsp.GemBonusRemain = info.CurrencyNum
		rsp.UpdateTime = mx.Unix(info.UpdateTime, 0)
	}
}

// GenGachaCost 生成抽卡成本
func GenGachaCost(s *enter.Session, cost *proto.ParcelCost) []*proto.ItemDB {
	itemLisl := make([]*proto.ItemDB, 0)
	for _, pi := range cost.ParcelInfos {
		switch pi.Key.Type {
		case proto.ParcelType_Item: // 物品
			game.RemoveItem(s, pi.Key.Id, int32(pi.Amount))
			itemLisl = append(itemLisl, game.GetItemDB(s, pi.Key.Id))
		case proto.ParcelType_Currency: // 代币
			game.UpCurrency(s, proto.CurrencyTypes(pi.Key.Id), -pi.Amount)
		}
	}

	return itemLisl
}
