package gdconf

import (
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadShopExcelTable() {
	g.GetExcel().ShopExcelTable = make([]*sro.ShopExcelTable, 0)
	name := "ShopExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().ShopExcelTable)
}

type ShopExcel struct {
	ShopExcelTableMap map[int64]*sro.ShopExcelTable
	ShopExcelTypeMap  map[string][]*sro.ShopExcelTable
}

func (g *GameConfig) gppShopExcelTable() {
	g.GetGPP().ShopExcel = &ShopExcel{
		ShopExcelTableMap: make(map[int64]*sro.ShopExcelTable),
		ShopExcelTypeMap:  make(map[string][]*sro.ShopExcelTable),
	}
	for _, v := range g.GetExcel().GetShopExcelTable() {
		g.GetGPP().ShopExcel.ShopExcelTableMap[v.Id] = v
		salePeriodTo, _ := time.Parse("2006-01-02 15:04:05", v.SalePeriodTo)
		if v.SalePeriodTo != "" && time.Now().After(salePeriodTo) {
			continue
		}
		if g.GetGPP().ShopExcel.ShopExcelTypeMap[v.CategoryType] == nil {
			g.GetGPP().ShopExcel.ShopExcelTypeMap[v.CategoryType] = make([]*sro.ShopExcelTable, 0)
		}
		g.GetGPP().ShopExcel.ShopExcelTypeMap[v.CategoryType] = append(
			g.GetGPP().ShopExcel.ShopExcelTypeMap[v.CategoryType], v)
	}
	logger.Info("处理商品配置完成,商店类型:%v个", len(g.GetGPP().ShopExcel.ShopExcelTypeMap))
}

func GetShopExcelType(categoryType string) []*sro.ShopExcelTable {
	if GC.GetGPP().ShopExcel.ShopExcelTypeMap == nil {
		return nil
	}
	return GC.GetGPP().ShopExcel.ShopExcelTypeMap[categoryType]
}

func GetShopExcelTable(shopId int64) *sro.ShopExcelTable {
	return GC.GetGPP().ShopExcel.ShopExcelTableMap[shopId]
}
