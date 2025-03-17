package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func ShopList(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.ShopListRequest)
	rsp := response.(*proto.ShopListResponse)

	rsp.ShopInfos = make([]*proto.ShopInfoDB, 0)
	rsp.ShopEligmaHistoryDBs = make([]*proto.ShopEligmaHistoryDB, 0)
	for _, categoryType := range req.CategoryList {
		conf := gdconf.GetShopInfoExcel(categoryType.String())
		info := &proto.ShopInfoDB{
			EventContentId:      0,
			Category:            categoryType,
			ManualRefreshCount:  0,                              // 手动刷新
			IsRefresh:           conf.IsRefresh,                 // 是否刷新
			NextAutoRefreshDate: time.Now().Add(24 * time.Hour), // 下一次
			LastAutoRefreshDate: time.Now(),                     // 上次刷新时间
			ShopProductList:     make([]*proto.ShopProductDB, 0),
		}
		if conf.IsRefresh {
			info.ShopProductList = game.GetRefreshShopProductList(categoryType)
		} else {
			info.ShopProductList = game.GetNoRefreshShopProductList(s, categoryType)
		}

		if len(info.ShopProductList) > 0 {
			rsp.ShopInfos = append(rsp.ShopInfos, info)
		}
	}
}

func ShopBuyRefreshMerchandise(s *enter.Session, request, response proto.Message) {
	// req := request.(*proto.ShopBuyRefreshMerchandiseRequest)
	// rsp := response.(*proto.ShopBuyRefreshMerchandiseResponse)

}

/*
0-20 1
20-40 2
40-60 3
60-80 4
5
*/

func ShopBuyEligma(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.ShopBuyEligmaRequest)
	rsp := response.(*proto.ShopBuyEligmaResponse)

	rsp.ParcelResultDB = game.ParcelResultDB(s, []*game.ParcelResult{
		{
			ParcelType: proto.ParcelType_Item,
			ParcelId:   req.CharacterUniqueId,
			Amount:     req.PurchaseCount,
		},
	})
	rsp.ConsumeResultDB = &proto.ConsumeResultDB{
		RemovedItemServerIds:                    make([]int64, 0),
		RemovedEquipmentServerIds:               make([]int64, 0),
		RemovedFurnitureServerIds:               make([]int64, 0),
		UsedItemServerIdAndRemainingCounts:      make(map[int64]int64),
		UsedEquipmentServerIdAndRemainingCounts: make(map[int64]int64),
		UsedFurnitureServerIdAndRemainingCounts: make(map[int64]int64),
	}
	conf := gdconf.GetShopExcelTable(req.ShopUniqueId)
	rsp.ShopProductDB = &proto.ShopProductDB{
		EventContentId:     0,
		ShopExcelId:        conf.GetId(),
		Category:           proto.ShopCategoryType_SecretStone,
		DisplayOrder:       conf.GetDisplayOrder(),
		PurchaseCount:      req.PurchaseCount,
		PurchaseCountLimit: conf.GetPurchaseCountLimit(),
		Price:              0,
		ProductType:        proto.ShopProductType_General,
	}
}

func ShopBuyMerchandise(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.ShopBuyMerchandiseRequest)
	rsp := response.(*proto.ShopBuyMerchandiseResponse)

	conf := gdconf.GetShopExcelTable(req.ShopUniqueId)
	if conf == nil {
		return
	}
	parcelResultList := make([]*game.ParcelResult, 0)
	// 购买物品
	for _, goodsId := range conf.GoodsId {
		goodsInfo := gdconf.GetGoodsExcelTable(goodsId)
		if goodsInfo == nil {
			continue
		}
		// 消耗
		parcelResultList = append(parcelResultList, game.GetParcelResultList(
			goodsInfo.ConsumeParcelType,
			goodsInfo.ConsumeParcelId,
			goodsInfo.ConsumeParcelAmount,
			true,
		)...)
		// 添加
		parcelResultList = append(parcelResultList, game.GetParcelResultList(
			goodsInfo.ParcelType_,
			goodsInfo.ParcelId,
			goodsInfo.ParcelAmount,
			false,
		)...)
	}
	// 构造回复
	rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
	rsp.ShopProductDB = &proto.ShopProductDB{
		EventContentId:     0,
		ShopExcelId:        conf.GetId(),
		Category:           proto.GetShopCategoryType(conf.GetCategoryType()),
		DisplayOrder:       conf.GetDisplayOrder(),
		PurchaseCount:      req.PurchaseCount,
		PurchaseCountLimit: conf.GetPurchaseCountLimit(),
		Price:              0,
		ProductType:        proto.ShopProductType_General,
	}
	rsp.ConsumeResultDB = &proto.ConsumeResultDB{
		RemovedItemServerIds:                    make([]int64, 0),
		RemovedEquipmentServerIds:               make([]int64, 0),
		RemovedFurnitureServerIds:               make([]int64, 0),
		UsedItemServerIdAndRemainingCounts:      make(map[int64]int64),
		UsedEquipmentServerIdAndRemainingCounts: make(map[int64]int64),
		UsedFurnitureServerIdAndRemainingCounts: make(map[int64]int64),
	}
}
