package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
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
	for id, _ := range addItemList {
		itemInfo := game.GetItemInfo(s, id)
		if itemInfo == nil {
			continue
		}
		rsp.AcquiredItems = append(rsp.AcquiredItems, &proto.ItemDB{
			Type:       proto.ParcelType_Item,
			ServerId:   itemInfo.ServerId,
			UniqueId:   itemInfo.UniqueId,
			StackCount: itemInfo.StackCount,
		})
	}
}

func ShopBuyGacha3(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopBuyGacha3Request)
	rsp := response.(*proto.ShopBuyGacha3Response)

	rsp.GachaResults = make([]*proto.GachaResult, 0)
	rsp.AcquiredItems = make([]*proto.ItemDB, 0)
	// 成本计算
	if game.UpCurrency(s, proto.CurrencyTypes_GemBonus, -req.Cost.Currency.Gem) != nil {
		num := req.Cost.Currency.Gem / 120

		isDu := false
		for _, itemDB := range req.Cost.ItemDBs {
			conf := gdconf.GetRecruitCoin()
			if conf == nil || conf.Id != itemDB.UniqueId {
				continue
			}
			itemInfo := game.GetItemInfo(s, itemDB.UniqueId)
			if itemInfo != nil &&
				itemDB.StackCount == 200 && itemInfo.UniqueId == itemDB.UniqueId &&
				game.RemoveItem(s, itemDB.UniqueId, itemDB.StackCount) {
				isDu = true
			}
		}

		var results []int64
		if isDu {
			results = game.GachaRun(1, true, false)
		} else {
			results = game.GachaRun(int(num), false, false)
		}

		list, addItemList := game.SaveGachaResults(s, results)
		rsp.GachaResults = list
		for id, _ := range addItemList {
			itemInfo := game.GetItemInfo(s, id)
			if itemInfo == nil {
				continue
			}
			rsp.AcquiredItems = append(rsp.AcquiredItems, &proto.ItemDB{
				Type:       proto.ParcelType_Item,
				ServerId:   itemInfo.ServerId,
				UniqueId:   itemInfo.UniqueId,
				StackCount: itemInfo.StackCount,
			})
		}
	}

	info := game.GetCurrencyInfo(s, proto.CurrencyTypes_GemBonus)
	if info != nil {
		rsp.GemBonusRemain = info.CurrencyNum
		rsp.UpdateTime = time.Unix(info.UpdateTime, 0)
	}
}
