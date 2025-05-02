package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadCafeRankExcelTable() {
	g.GetExcel().CafeRankExcelTable = make([]*sro.CafeRankExcelTable, 0)
	name := "CafeRankExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().CafeRankExcelTable)
}

type CafeRankExcel struct {
	CafeRankExcelMap map[int64]map[int32]*sro.CafeRankExcelTable
}

func (g *GameConfig) gppCafeRankExcelTable() {
	g.GetGPP().CafeRankExcel = &CafeRankExcel{
		CafeRankExcelMap: make(map[int64]map[int32]*sro.CafeRankExcelTable),
	}

	for _, v := range g.GetExcel().GetCafeRankExcelTable() {
		if g.GetGPP().CafeRankExcel.CafeRankExcelMap[v.CafeId] == nil {
			g.GetGPP().CafeRankExcel.CafeRankExcelMap[v.CafeId] = make(map[int32]*sro.CafeRankExcelTable)
		}
		g.GetGPP().CafeRankExcel.CafeRankExcelMap[v.CafeId][v.Rank] = v
	}

	logger.Info("处理咖啡厅等级配置表完成,数量:%v个",
		len(g.GetGPP().CafeRankExcel.CafeRankExcelMap))
}

func GetCafeRankExcelTable(cafeId int64, rank int32) *sro.CafeRankExcelTable {
	if GC.GetGPP().CafeRankExcel.CafeRankExcelMap[cafeId] == nil {
		return nil
	}
	return GC.GetGPP().CafeRankExcel.CafeRankExcelMap[cafeId][rank]
}
