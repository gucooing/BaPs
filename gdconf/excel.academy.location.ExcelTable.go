package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyLocationExcelTable() {
	g.GetExcel().AcademyLocationExcelTable = make([]*sro.AcademyLocationExcelTable, 0)
	name := "AcademyLocationExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().AcademyLocationExcelTable)
}

type AcademyLocationExcel struct {
	AcademyLocationExcelMap map[int64]*sro.AcademyLocationExcelTable
}

func (g *GameConfig) gppAcademyLocationExcelTable() {
	g.GetGPP().AcademyLocationExcel = &AcademyLocationExcel{
		AcademyLocationExcelMap: make(map[int64]*sro.AcademyLocationExcelTable),
	}
	for _, v := range g.GetExcel().GetAcademyLocationExcelTable() {
		g.GetGPP().AcademyLocationExcel.AcademyLocationExcelMap[v.Id] = v
	}

	logger.Info("处理课程表学院信息完成,数量:%v个",
		len(g.GetGPP().AcademyLocationExcel.AcademyLocationExcelMap))
}

func GetAcademyLocationExcelTableList() []*sro.AcademyLocationExcelTable {
	return GC.GetExcel().GetAcademyLocationExcelTable()
}
