package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadTimeAttackDungeonRewardExcelTable() {
	g.GetExcel().TimeAttackDungeonRewardExcelTable = make([]*sro.TimeAttackDungeonRewardExcelTable, 0)
	name := "TimeAttackDungeonRewardExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().TimeAttackDungeonRewardExcelTable)
}

type TimeAttackDungeonRewardExcel struct {
	TimeAttackDungeonRewardExcelMap map[int64]*sro.TimeAttackDungeonRewardExcelTable
}

func (g *GameConfig) gppTimeAttackDungeonRewardExcelTable() {
	g.GetGPP().TimeAttackDungeonRewardExcel = &TimeAttackDungeonRewardExcel{
		TimeAttackDungeonRewardExcelMap: make(map[int64]*sro.TimeAttackDungeonRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetTimeAttackDungeonRewardExcelTable() {
		g.GetGPP().TimeAttackDungeonRewardExcel.TimeAttackDungeonRewardExcelMap[v.Id] = v
	}

	logger.Info("处理综合战术考试关卡配置完成,技能配置:%v个",
		len(g.GetGPP().TimeAttackDungeonRewardExcel.TimeAttackDungeonRewardExcelMap))
}

func GetTimeAttackDungeonRewardExcelTable(id int64) *sro.TimeAttackDungeonRewardExcelTable {
	return GC.GetGPP().TimeAttackDungeonRewardExcel.TimeAttackDungeonRewardExcelMap[id]
}
