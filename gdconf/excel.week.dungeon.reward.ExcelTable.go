package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadWeekDungeonRewardExcelTable() {
	g.GetExcel().WeekDungeonRewardExcelTable = make([]*sro.WeekDungeonRewardExcelTable, 0)
	name := "WeekDungeonRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().WeekDungeonRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetWeekDungeonRewardExcelTable()))
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
