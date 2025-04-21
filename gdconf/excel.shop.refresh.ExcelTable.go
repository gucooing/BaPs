package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadShopRefreshExcelTable() {
	g.GetExcel().ShopRefreshExcelTable = make([]*sro.ShopRefreshExcelTable, 0)
	name := "ShopRefreshExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().ShopRefreshExcelTable)
}

type ShopRefreshExcel struct {
	ShopRefreshExcelTableMap  map[string][]*sro.ShopRefreshExcelTable
	ShopRefreshExcelTableList map[int64]*sro.ShopRefreshExcelTable
}

func (g *GameConfig) gppShopRefreshExcelTable() {
	g.GetGPP().ShopRefreshExcel = &ShopRefreshExcel{
		ShopRefreshExcelTableMap:  make(map[string][]*sro.ShopRefreshExcelTable),
		ShopRefreshExcelTableList: make(map[int64]*sro.ShopRefreshExcelTable),
	}
	for _, v := range g.GetExcel().GetShopRefreshExcelTable() {
		g.GetGPP().ShopRefreshExcel.ShopRefreshExcelTableList[v.Id] = v
		if g.GetGPP().ShopRefreshExcel.ShopRefreshExcelTableMap[v.CategoryType] == nil {
			g.GetGPP().ShopRefreshExcel.ShopRefreshExcelTableMap[v.CategoryType] = make([]*sro.ShopRefreshExcelTable, 0)
		}
		g.GetGPP().ShopRefreshExcel.ShopRefreshExcelTableMap[v.CategoryType] = append(
			g.GetGPP().ShopRefreshExcel.ShopRefreshExcelTableMap[v.CategoryType],
			v,
		)
	}
	logger.Info("处理可刷新商品配置完成,商店类型:%v个", len(g.GetGPP().ShopRefreshExcel.ShopRefreshExcelTableMap))
}

func GetShopRefreshExcelMap(categoryType string) []*sro.ShopRefreshExcelTable {
	return GC.GetGPP().ShopRefreshExcel.ShopRefreshExcelTableMap[categoryType]
}

func GetShopRefreshExcel(shopId int64) *sro.ShopRefreshExcelTable {
	return GC.GetGPP().ShopRefreshExcel.ShopRefreshExcelTableList[shopId]
}
