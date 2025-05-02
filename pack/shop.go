package pack

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func ShopList(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopListRequest)
	rsp := response.(*proto.ShopListResponse)

	rsp.ShopInfos = make([]*proto.ShopInfoDB, 0)
	rsp.ShopEligmaHistoryDBs = make([]*proto.ShopEligmaHistoryDB, 0)
	for _, categoryType := range req.CategoryList {
		conf := gdconf.GetShopInfoExcel(categoryType.String())
		info := &proto.ShopInfoDB{
			EventContentId:      0,
			Category:            categoryType,
			ManualRefreshCount:  0,                            // 手动刷新
			IsRefresh:           conf.IsRefresh,               // 是否刷新
			NextAutoRefreshDate: mx.Now().Add(24 * time.Hour), // 下一次
			LastAutoRefreshDate: mx.Now(),                     // 上次刷新时间
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

func ShopBuyRefreshMerchandise(s *enter.Session, request, response mx.Message) {
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

func ShopBuyEligma(s *enter.Session, request, response mx.Message) {
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

func ShopBuyMerchandise(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopBuyMerchandiseRequest)
	rsp := response.(*proto.ShopBuyMerchandiseResponse)

	parcelResultList := make([]*game.ParcelResult, 0)
	addParcelResult := func(typeList []string, idList, numList []int64, purchaseCount int64, isDel bool) {
		//  不验有没有那么多了,随意了
		if len(typeList) == len(idList) &&
			len(idList) == len(numList) {
			for index, rewardType := range typeList {
				num := numList[index] * purchaseCount
				if isDel {
					num = -numList[index]
				}
				parcelResultList = append(parcelResultList, &game.ParcelResult{
					ParcelType: proto.ParcelType_None.Value(rewardType),
					ParcelId:   idList[index],
					Amount:     num,
				})
			}
		}
	}
	addGoodsExcel := func(goodsInfo *sro.GoodsExcelTable) {
		if goodsInfo == nil {
			return
		}
		// 消耗
		addParcelResult(goodsInfo.ConsumeParcelType,
			goodsInfo.ConsumeParcelId,
			goodsInfo.ConsumeParcelAmount,
			req.PurchaseCount,
			true)
		// 添加
		addParcelResult(
			goodsInfo.ParcelType,
			goodsInfo.ParcelId,
			goodsInfo.ParcelAmount,
			req.PurchaseCount,
			false,
		)
	}

	// 购买物品
	var productType proto.ShopProductType
	var category proto.ShopCategoryType
	var purchaseCountLimit int64
	var displayOrder int64
	if req.IsRefreshGoods {
		conf := gdconf.GetShopRefreshExcel(req.ShopUniqueId)
		productType = proto.ShopProductType_Refresh
		category = proto.ShopCategoryType_General.Value(conf.GetCategoryType())
		purchaseCountLimit = 1
		displayOrder = conf.GetDisplayOrder()

		addGoodsExcel(gdconf.GetGoodsExcelTable(conf.GoodsId))
	} else {
		conf := gdconf.GetShopExcelTable(req.ShopUniqueId)
		productType = proto.ShopProductType_General
		category = proto.ShopCategoryType_General.Value(conf.GetCategoryType())
		purchaseCountLimit = conf.GetPurchaseCountLimit()
		displayOrder = conf.GetDisplayOrder()

		for _, goodsId := range conf.GoodsId {
			addGoodsExcel(gdconf.GetGoodsExcelTable(goodsId))
		}
	}

	// 构造回复
	rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
	rsp.ShopProductDB = &proto.ShopProductDB{
		EventContentId:     0,
		ShopExcelId:        req.ShopUniqueId,
		Category:           category,
		DisplayOrder:       displayOrder,
		PurchaseCount:      req.PurchaseCount,
		PurchaseCountLimit: purchaseCountLimit,
		Price:              0,
		ProductType:        productType,
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
