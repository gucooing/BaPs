package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadGuideMissionExcelTable() {
	g.GetExcel().GuideMissionExcelTable = make([]*sro.GuideMissionExcelTable, 0)
	name := "GuideMissionExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().GuideMissionExcelTable)
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
