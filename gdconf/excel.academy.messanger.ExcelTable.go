package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyMessangerExcelTable() {
	g.GetExcel().AcademyMessangerExcelTable = make([]*sro.AcademyMessangerExcelTable, 0)
	nameList := []string{"AcademyMessanger1ExcelTable.json", "AcademyMessanger2ExcelTable.json",
		"AcademyMessanger3ExcelTable.json", "AcademyMessangerExcelTable.json"}
	for _, name := range nameList {
		list := make([]*sro.AcademyMessangerExcelTable, 0)
		loadExcelJson(g.excelPath+name, &list)
		g.GetExcel().AcademyMessangerExcelTable = append(g.GetExcel().AcademyMessangerExcelTable, list...)
	}
}

type AcademyMessangerExcel struct {
	AcademyMessangerExcelMap map[int64]*sro.AcademyMessangerExcelTable
}

func (g *GameConfig) gppAcademyMessangerExcelTable() {
	g.GetGPP().AcademyMessangerExcel = &AcademyMessangerExcel{
		AcademyMessangerExcelMap: make(map[int64]*sro.AcademyMessangerExcelTable),
	}

	for _, v := range g.GetExcel().GetAcademyMessangerExcelTable() {
		g.GetGPP().AcademyMessangerExcel.AcademyMessangerExcelMap[v.MessageGroupId] = v
	}

	logger.Info("处理MomoTalk对话配置完成,MomoTalk对话:%v个",
		len(g.GetGPP().AcademyMessangerExcel.AcademyMessangerExcelMap))
}

func GetAcademyMessangerExcelTable(gid int64) *sro.AcademyMessangerExcelTable {
	return GC.GetGPP().AcademyMessangerExcel.AcademyMessangerExcelMap[gid]
}
