package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadRaidStageSeasonRewardExcelTable() {
	g.GetExcel().RaidStageSeasonRewardExcelTable = make([]*sro.RaidStageSeasonRewardExcelTable, 0)
	name := "RaidStageSeasonRewardExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().RaidStageSeasonRewardExcelTable)
}

type RaidStageSeasonRewardExcel struct {
	RaidStageSeasonRewardExcelMap map[int64]*sro.RaidStageSeasonRewardExcelTable
}

func (g *GameConfig) gppRaidStageSeasonRewardExcelTable() {
	g.GetGPP().RaidStageSeasonRewardExcel = &RaidStageSeasonRewardExcel{
		RaidStageSeasonRewardExcelMap: make(map[int64]*sro.RaidStageSeasonRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetRaidStageSeasonRewardExcelTable() {
		g.GetGPP().RaidStageSeasonRewardExcel.RaidStageSeasonRewardExcelMap[v.SeasonRewardId] = v
	}

	logger.Info("处理总力战总分奖励配置表完成,总力战总分奖励配置:%v个",
		len(g.GetGPP().RaidStageSeasonRewardExcel.RaidStageSeasonRewardExcelMap))
}

func GetRaidStageSeasonRewardExcelTable(id int64) *sro.RaidStageSeasonRewardExcelTable {
	return GC.GetGPP().RaidStageSeasonRewardExcel.RaidStageSeasonRewardExcelMap[id]
}
