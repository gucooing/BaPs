package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEliminateRaidStageSeasonRewardExcelTable() {
	g.GetExcel().EliminateRaidStageSeasonRewardExcelTable = make([]*sro.EliminateRaidStageSeasonRewardExcelTable, 0)
	name := "EliminateRaidStageSeasonRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EliminateRaidStageSeasonRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetEliminateRaidStageSeasonRewardExcelTable()))
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
