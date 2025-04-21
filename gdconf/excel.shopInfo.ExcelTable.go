package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadShopInfoExcelTable() {
	g.GetExcel().ShopInfoExcelTable = make([]*sro.ShopInfoExcelTable, 0)
	name := "ShopInfoExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().ShopInfoExcelTable)
}

type ShopInfoExcel struct {
	ShopInfoExcelTableMap map[string]*sro.ShopInfoExcelTable
}

func (g *GameConfig) gppShopInfoExcelTable() {
	g.GetGPP().ShopInfoExcel = &ShopInfoExcel{
		ShopInfoExcelTableMap: make(map[string]*sro.ShopInfoExcelTable),
	}
	for _, v := range g.GetExcel().GetShopInfoExcelTable() {
		g.GetGPP().ShopInfoExcel.ShopInfoExcelTableMap[v.CategoryType] = v
	}
	logger.Info("处理商店配置完成,商店类型:%v个", len(g.GetGPP().ShopInfoExcel.ShopInfoExcelTableMap))
}

func GetShopInfoExcel(categoryType string) *sro.ShopInfoExcelTable {
	return GC.GetGPP().ShopInfoExcel.ShopInfoExcelTableMap[categoryType]
}
