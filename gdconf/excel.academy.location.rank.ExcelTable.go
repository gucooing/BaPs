package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadAcademyLocationRankExcelTable() {
	g.GetExcel().AcademyLocationRankExcelTable = make([]*sro.AcademyLocationRankExcelTable, 0)
	name := "AcademyLocationRankExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().AcademyLocationRankExcelTable)
}

type AcademyLocationRankExcel struct {
	AcademyLocationRankExcelMap map[int64]*sro.AcademyLocationRankExcelTable
}

func (g *GameConfig) gppAcademyLocationRankExcelTable() {
	g.GetGPP().AcademyLocationRankExcel = &AcademyLocationRankExcel{
		AcademyLocationRankExcelMap: make(map[int64]*sro.AcademyLocationRankExcelTable),
	}
	for _, v := range g.GetExcel().GetAcademyLocationRankExcelTable() {
		g.GetGPP().AcademyLocationRankExcel.AcademyLocationRankExcelMap[v.Rank] = v
	}

	logger.Info("处理课程表等级配置完成,数量:%v个",
		len(g.GetGPP().AcademyLocationRankExcel.AcademyLocationRankExcelMap))
}

func GetAcademyLocationRankExcelTable(rank int64) *sro.AcademyLocationRankExcelTable {
	return GC.GetGPP().AcademyLocationRankExcel.AcademyLocationRankExcelMap[rank]
}
