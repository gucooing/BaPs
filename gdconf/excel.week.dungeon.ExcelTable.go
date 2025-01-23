package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadWeekDungeonExcelTable() {
	g.GetExcel().WeekDungeonExcelTable = make([]*sro.WeekDungeonExcelTable, 0)
	name := "WeekDungeonExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().WeekDungeonExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetWeekDungeonExcelTable()))
}

type WeekDungeonExcel struct {
	WeekDungeonExcelMap map[int64]*sro.WeekDungeonExcelTable
}

func (g *GameConfig) gppWeekDungeonExcelTable() {
	g.GetGPP().WeekDungeonExcel = &WeekDungeonExcel{
		WeekDungeonExcelMap: make(map[int64]*sro.WeekDungeonExcelTable, 0),
	}

	for _, v := range g.GetExcel().GetWeekDungeonExcelTable() {
		g.GetGPP().WeekDungeonExcel.WeekDungeonExcelMap[v.StageId] = v
	}

	logger.Info("悬赏通缉关卡信息数量完成:%v个", len(g.GetGPP().WeekDungeonExcel.WeekDungeonExcelMap))
}

func GetWeekDungeonExcelTable(stageId int64) *sro.WeekDungeonExcelTable {
	return GC.GetGPP().WeekDungeonExcel.WeekDungeonExcelMap[stageId]
}
