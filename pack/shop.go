package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func ShopList(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopListRequest)
	rsp := response.(*proto.ShopListResponse)

	rsp.ShopInfos = make([]*proto.ShopInfoDB, 0)
	rsp.ShopEligmaHistoryDBs = make([]*proto.ShopEligmaHistoryDB, 0)

	for _, categoryType := range req.CategoryList {
		if categoryType == proto.ShopCategoryType_General {
			continue
		}
		conf := gdconf.GetShopInfoExcel(categoryType.String())
		info := &proto.ShopInfoDB{
			EventContentId:      0,
			Category:            categoryType,
			ManualRefreshCount:  0,
			IsRefresh:           conf.IsRefresh,
			NextAutoRefreshDate: time.Now().Add(24 * time.Hour),
			LastAutoRefreshDate: time.Now(),
			ShopProductList:     make([]*proto.ShopProductDB, 0),
		}
		for _, product := range gdconf.GetShopExcelType(categoryType.String()) {
			info.ShopProductList = append(info.ShopProductList, &proto.ShopProductDB{
				EventContentId:     0,
				ShopExcelId:        product.Id,
				Category:           proto.ShopCategoryType(proto.ShopCategoryType_value[product.CategoryType]),
				DisplayOrder:       product.DisplayOrder,
				PurchaseCount:      0,
				PurchaseCountLimit: product.PurchaseCountLimit,
				Price:              0,
				ProductType:        proto.ShopProductType_General,
			})
		}

		if len(info.ShopProductList) > 0 {
			rsp.ShopInfos = append(rsp.ShopInfos, info)
		}
	}
}
