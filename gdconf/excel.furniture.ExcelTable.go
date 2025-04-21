package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadFurnitureExcelTable() {
	g.GetExcel().FurnitureExcelTable = make([]*sro.FurnitureExcelTable, 0)
	name := "FurnitureExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().FurnitureExcelTable)
}

type FurnitureExcel struct {
	FurnitureExcelTableMap map[int64]*sro.FurnitureExcelTable
}

func (g *GameConfig) gppFurnitureExcelTable() {
	g.GetGPP().FurnitureExcel = &FurnitureExcel{
		FurnitureExcelTableMap: make(map[int64]*sro.FurnitureExcelTable, 0),
	}
	for _, v := range g.GetExcel().GetFurnitureExcelTable() {
		g.GetGPP().FurnitureExcel.FurnitureExcelTableMap[v.Id] = v
	}
	logger.Info("处理家具配置表完成数量:%v个", len(g.GetGPP().FurnitureExcel.FurnitureExcelTableMap))
}

func GetFurnitureExcelTable(id int64) *sro.FurnitureExcelTable {
	return GC.GetGPP().FurnitureExcel.FurnitureExcelTableMap[id]
}

func GetFurnitureExcelTableMap() []*sro.FurnitureExcelTable {
	return GC.GetExcel().GetFurnitureExcelTable()
}
