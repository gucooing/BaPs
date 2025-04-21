package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadEliminateRaidStageRewardExcelTable() {
	g.GetExcel().EliminateRaidStageRewardExcelTable = make([]*sro.EliminateRaidStageRewardExcelTable, 0)
	name := "EliminateRaidStageRewardExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().EliminateRaidStageRewardExcelTable)
}

type EliminateRaidStageRewardExcel struct {
	EliminateRaidStageRewardExcelMap map[int64][]*sro.EliminateRaidStageRewardExcelTable
}

func (g *GameConfig) gppEliminateRaidStageRewardExcelTable() {
	g.GetGPP().EliminateRaidStageRewardExcel = &EliminateRaidStageRewardExcel{
		EliminateRaidStageRewardExcelMap: make(map[int64][]*sro.EliminateRaidStageRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetEliminateRaidStageRewardExcelTable() {
		if g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[v.GroupId] == nil {
			g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[v.GroupId] = make([]*sro.EliminateRaidStageRewardExcelTable, 0)
		}
		g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[v.GroupId] =
			append(g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[v.GroupId], v)
	}

	logger.Info("处理大决战关卡通过奖励配置表完成,总力战关卡通过奖励配置:%v个",
		len(g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap))
}

func GetEliminateRaidStageRewardExcelTable(gid int64) []*sro.EliminateRaidStageRewardExcelTable {
	return GC.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[gid]
}
