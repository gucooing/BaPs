package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadRaidStageRewardExcelTable() {
	g.GetExcel().RaidStageRewardExcelTable = make([]*sro.RaidStageRewardExcelTable, 0)
	name := "RaidStageRewardExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().RaidStageRewardExcelTable)
}

type RaidStageRewardExcel struct {
	RaidStageRewardExcelMap map[int64][]*sro.RaidStageRewardExcelTable
}

func (g *GameConfig) gppRaidStageRewardExcelTable() {
	g.GetGPP().RaidStageRewardExcel = &RaidStageRewardExcel{
		RaidStageRewardExcelMap: make(map[int64][]*sro.RaidStageRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetRaidStageRewardExcelTable() {
		if g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[v.GroupId] == nil {
			g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[v.GroupId] = make([]*sro.RaidStageRewardExcelTable, 0)
		}
		g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[v.GroupId] =
			append(g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[v.GroupId], v)
	}

	logger.Info("处理总力战关卡通过奖励配置表完成,总力战关卡通过奖励配置:%v个",
		len(g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap))
}

func GetRaidStageRewardExcelTable(gid int64) []*sro.RaidStageRewardExcelTable {
	return GC.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[gid]
}
