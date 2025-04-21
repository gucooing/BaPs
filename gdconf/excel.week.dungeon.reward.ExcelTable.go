package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadWeekDungeonRewardExcelTable() {
	g.GetExcel().WeekDungeonRewardExcelTable = make([]*sro.WeekDungeonRewardExcelTable, 0)
	name := "WeekDungeonRewardExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().WeekDungeonRewardExcelTable)
}

type WeekDungeonRewardExcel struct {
	WeekDungeonRewardExcelMap map[int64][]*sro.WeekDungeonRewardExcelTable
}

func (g *GameConfig) gppWeekDungeonRewardExcelTable() {
	g.GetGPP().WeekDungeonRewardExcel = &WeekDungeonRewardExcel{
		WeekDungeonRewardExcelMap: make(map[int64][]*sro.WeekDungeonRewardExcelTable, 0),
	}

	for _, v := range g.GetExcel().GetWeekDungeonRewardExcelTable() {
		if g.GetGPP().WeekDungeonRewardExcel.WeekDungeonRewardExcelMap[v.GroupId] == nil {
			g.GetGPP().WeekDungeonRewardExcel.WeekDungeonRewardExcelMap[v.GroupId] = make([]*sro.WeekDungeonRewardExcelTable, 0)
		}
		g.GetGPP().WeekDungeonRewardExcel.WeekDungeonRewardExcelMap[v.GroupId] = append(
			g.GetGPP().WeekDungeonRewardExcel.WeekDungeonRewardExcelMap[v.GroupId], v)
	}

	logger.Info("处理悬赏通缉关卡奖励配置完成,数量:%v个", len(g.GetGPP().WeekDungeonRewardExcel.WeekDungeonRewardExcelMap))
}

func GetWeekDungeonRewardExcelList(stageId int64) []*sro.WeekDungeonRewardExcelTable {
	return GC.GetGPP().WeekDungeonRewardExcel.WeekDungeonRewardExcelMap[stageId]
}
