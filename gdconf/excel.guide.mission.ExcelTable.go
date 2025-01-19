package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadGuideMissionExcelTable() {
	g.GetExcel().GuideMissionExcelTable = make([]*sro.GuideMissionExcelTable, 0)
	name := "GuideMissionExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().GuideMissionExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GuideMissionExcelTable))
}

type GuideMissionExcel struct {
	GuideMissionExcelMap map[int64]*sro.GuideMissionExcelTable
}

func (g *GameConfig) gppGuideMissionExcelTable() {
	g.GetGPP().GuideMissionExcel = &GuideMissionExcel{
		GuideMissionExcelMap: make(map[int64]*sro.GuideMissionExcelTable),
	}
	for _, v := range g.GetExcel().GetGuideMissionExcelTable() {
		g.GetGPP().GuideMissionExcel.GuideMissionExcelMap[v.Id] = v
	}

	logger.Info("处理成就配置完成,成就:%v个",
		len(g.GetGPP().GuideMissionExcel.GuideMissionExcelMap))
}

func GetGuideMissionExcelTable(id int64) *sro.GuideMissionExcelTable {
	return GC.GetGPP().GuideMissionExcel.GuideMissionExcelMap[id]
}
