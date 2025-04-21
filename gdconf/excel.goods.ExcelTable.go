package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadGoodsExcelTable() {
	g.GetExcel().GoodsExcelTable = make([]*sro.GoodsExcelTable, 0)
	name := "GoodsExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().GoodsExcelTable)
}

type GoodsExcel struct {
	GoodsExcelMap map[int64]*sro.GoodsExcelTable
}

func (g *GameConfig) gppGoodsExcelTable() {
	g.GetGPP().GoodsExcel = &GoodsExcel{
		GoodsExcelMap: make(map[int64]*sro.GoodsExcelTable),
	}
	for _, v := range g.GetExcel().GetGoodsExcelTable() {
		g.GetGPP().GoodsExcel.GoodsExcelMap[v.Id] = v
	}

	logger.Info("处理商品配置完成,成就:%v个",
		len(g.GetGPP().GoodsExcel.GoodsExcelMap))
}

func GetGoodsExcelTable(id int64) *sro.GoodsExcelTable {
	return GC.GetGPP().GoodsExcel.GoodsExcelMap[id]
}
