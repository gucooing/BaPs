package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func ShopList(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ShopListRequest)
	rsp := response.(*proto.ShopListResponse)

	rsp.ShopInfos = make([]*proto.ShopInfoDB, 0)
	rsp.ShopEligmaHistoryDBs = make([]*proto.ShopEligmaHistoryDB, 0)

	return

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
			info.ShopProductList = game.GetRefreshShopProductList(categoryType.String())
		} else {
			info.ShopProductList = game.GetNoRefreshShopProductList(s, categoryType.String())
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

}
