package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadRaidStageSeasonRewardExcelTable() {
	g.GetExcel().RaidStageSeasonRewardExcelTable = make([]*sro.RaidStageSeasonRewardExcelTable, 0)
	name := "RaidStageSeasonRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().RaidStageSeasonRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetRaidStageSeasonRewardExcelTable()))
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
