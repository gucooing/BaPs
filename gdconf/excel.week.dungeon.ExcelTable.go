package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadWeekDungeonExcelTable() {
	g.GetExcel().WeekDungeonExcelTable = make([]*sro.WeekDungeonExcelTable, 0)
	name := "WeekDungeonExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().WeekDungeonExcelTable)
}

type WeekDungeonExcel struct {
	WeekDungeonExcelMap map[int64]*sro.WeekDungeonExcelTable
}

func (g *GameConfig) gppWeekDungeonExcelTable() {
	g.GetGPP().WeekDungeonExcel = &WeekDungeonExcel{
		WeekDungeonExcelMap: make(map[int64]*sro.WeekDungeonExcelTable, 0),
	}

	for _, v := range g.GetExcel().GetWeekDungeonExcelTable() {
		if len(v.StarGoal) != len(v.StarGoalAmount) {
			logger.Warn("WeekDungeonExcelTable.json StarGoal和StarGoalAmount不对应")
		}
		g.GetGPP().WeekDungeonExcel.WeekDungeonExcelMap[v.StageId] = v
	}

	logger.Info("悬赏通缉关卡信息数量完成:%v个", len(g.GetGPP().WeekDungeonExcel.WeekDungeonExcelMap))
}

func GetWeekDungeonExcelTable(stageId int64) *sro.WeekDungeonExcelTable {
	return GC.GetGPP().WeekDungeonExcel.WeekDungeonExcelMap[stageId]
}
