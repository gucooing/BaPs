package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadEliminateRaidStageSeasonRewardExcelTable() {
	g.GetExcel().EliminateRaidStageSeasonRewardExcelTable = make([]*sro.EliminateRaidStageSeasonRewardExcelTable, 0)
	name := "EliminateRaidStageSeasonRewardExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().EliminateRaidStageSeasonRewardExcelTable)
}

type EliminateRaidStageSeasonRewardExcel struct {
	EliminateRaidStageSeasonRewardExcelMap map[int64]*sro.EliminateRaidStageSeasonRewardExcelTable
}

func (g *GameConfig) gppEliminateRaidStageSeasonRewardExcelTable() {
	g.GetGPP().EliminateRaidStageSeasonRewardExcel = &EliminateRaidStageSeasonRewardExcel{
		EliminateRaidStageSeasonRewardExcelMap: make(map[int64]*sro.EliminateRaidStageSeasonRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetEliminateRaidStageSeasonRewardExcelTable() {
		g.GetGPP().EliminateRaidStageSeasonRewardExcel.EliminateRaidStageSeasonRewardExcelMap[v.SeasonRewardId] = v
	}

	logger.Info("处理大决战总分奖励配置表完成,大决战总分奖励配置:%v个",
		len(g.GetGPP().EliminateRaidStageSeasonRewardExcel.EliminateRaidStageSeasonRewardExcelMap))
}

func GetEliminateRaidStageSeasonRewardExcelTable(id int64) *sro.EliminateRaidStageSeasonRewardExcelTable {
	return GC.GetGPP().EliminateRaidStageSeasonRewardExcel.EliminateRaidStageSeasonRewardExcelMap[id]
}
