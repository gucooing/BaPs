package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadEliminateRaidStageExcelTable() {
	g.GetExcel().EliminateRaidStageExcelTable = make([]*sro.EliminateRaidStageExcelTable, 0)
	name := "EliminateRaidStageExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().EliminateRaidStageExcelTable)
}

type EliminateRaidStageExcel struct {
	EliminateRaidStageExcelMap map[int64]*sro.EliminateRaidStageExcelTable
}

func (g *GameConfig) gppEliminateRaidStageExcelTable() {
	g.GetGPP().EliminateRaidStageExcel = &EliminateRaidStageExcel{
		EliminateRaidStageExcelMap: make(map[int64]*sro.EliminateRaidStageExcelTable),
	}

	for _, v := range g.GetExcel().GetEliminateRaidStageExcelTable() {
		g.GetGPP().EliminateRaidStageExcel.EliminateRaidStageExcelMap[v.Id] = v
	}

	logger.Info("处理大决战关卡配置表完成,关卡配置:%v个",
		len(g.GetGPP().EliminateRaidStageExcel.EliminateRaidStageExcelMap))
}

func GetEliminateRaidStageExcelTable(id int64) *sro.EliminateRaidStageExcelTable {
	return GC.GetGPP().EliminateRaidStageExcel.EliminateRaidStageExcelMap[id]
}
