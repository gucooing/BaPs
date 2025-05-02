package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadTimeAttackDungeonGeasExcelTable() {
	g.GetExcel().TimeAttackDungeonGeasExcelTable = make([]*sro.TimeAttackDungeonGeasExcelTable, 0)
	name := "TimeAttackDungeonGeasExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().TimeAttackDungeonGeasExcelTable)
}

type TimeAttackDungeonGeasExcel struct {
	TimeAttackDungeonGeasExcelMap map[int64]*sro.TimeAttackDungeonGeasExcelTable
}

func (g *GameConfig) gppTimeAttackDungeonGeasExcelTable() {
	g.GetGPP().TimeAttackDungeonGeasExcel = &TimeAttackDungeonGeasExcel{
		TimeAttackDungeonGeasExcelMap: make(map[int64]*sro.TimeAttackDungeonGeasExcelTable),
	}

	for _, v := range g.GetExcel().GetTimeAttackDungeonGeasExcelTable() {
		g.GetGPP().TimeAttackDungeonGeasExcel.TimeAttackDungeonGeasExcelMap[v.Id] = v
	}

	logger.Info("处理综合战术考试关卡配置完成,技能配置:%v个",
		len(g.GetGPP().TimeAttackDungeonGeasExcel.TimeAttackDungeonGeasExcelMap))
}

func GetTimeAttackDungeonGeasExcelTable(id int64) *sro.TimeAttackDungeonGeasExcelTable {
	return GC.GetGPP().TimeAttackDungeonGeasExcel.TimeAttackDungeonGeasExcelMap[id]
}
