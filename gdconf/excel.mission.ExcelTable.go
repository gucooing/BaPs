package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadMissionExcelTable() {
	g.GetExcel().MissionExcelTable = make([]*sro.MissionExcelTable, 0)
	name := "MissionExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().MissionExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetMissionExcelTable()))
}

type MissionExcel struct {
	MissionExcelTableMap      map[int64]*sro.MissionExcelTable
	MissionExcelTableCategory map[string][]*sro.MissionExcelTable
}

func (g *GameConfig) gppMissionExcelTable() {
	g.GetGPP().MissionExcel = &MissionExcel{
		MissionExcelTableMap:      make(map[int64]*sro.MissionExcelTable, 0),
		MissionExcelTableCategory: make(map[string][]*sro.MissionExcelTable),
	}
	for _, v := range g.GetExcel().GetMissionExcelTable() {
		g.GetGPP().MissionExcel.MissionExcelTableMap[v.Id] = v
		if g.GetGPP().MissionExcel.MissionExcelTableCategory[v.Category] == nil {
			g.GetGPP().MissionExcel.MissionExcelTableCategory[v.Category] = make([]*sro.MissionExcelTable, 0)
		}
		g.GetGPP().MissionExcel.MissionExcelTableCategory[v.Category] =
			append(g.GetGPP().MissionExcel.MissionExcelTableCategory[v.Category], v)
	}
	logger.Info("处理任务配置表完成数量:%v个", len(g.GetGPP().MissionExcel.MissionExcelTableMap))
}

func GetMissionExcelTableCategoryList(category string) []*sro.MissionExcelTable {
	return GC.GetGPP().MissionExcel.MissionExcelTableCategory[category]
}

func GetMissionExcelTable(id int64) *sro.MissionExcelTable {
	return GC.GetGPP().MissionExcel.MissionExcelTableMap[id]
}
