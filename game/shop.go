package game

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetNoRefreshShopProductList(s *enter.Session, categoryType string) []*proto.ShopProductDB {
	list := make([]*proto.ShopProductDB, 0)
	if categoryType == "SecretStone" {
		return list
	}
	for _, product := range gdconf.GetShopExcelType(categoryType) {
		list = append(list, &proto.ShopProductDB{
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

	return list
}

func GetRefreshShopProductList(categoryType string) []*proto.ShopProductDB {
	list := make([]*proto.ShopProductDB, 0)
	for _, product := range gdconf.GetShopExcelType(categoryType) {
		list = append(list, &proto.ShopProductDB{
			EventContentId:     0,
			ShopExcelId:        product.Id,
			Category:           proto.ShopCategoryType(proto.ShopCategoryType_value[product.CategoryType]),
			DisplayOrder:       product.DisplayOrder,
			PurchaseCount:      0,
			PurchaseCountLimit: product.PurchaseCountLimit,
			Price:              0,
			ProductType:        proto.ShopProductType_Refresh,
		})
	}

	return list
}
