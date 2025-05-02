package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadAcademyFavorScheduleExcelTable() {
	g.GetExcel().AcademyFavorScheduleExcelTable = make([]*sro.AcademyFavorScheduleExcelTable, 0)
	name := "AcademyFavorScheduleExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().AcademyFavorScheduleExcelTable)
}

type AcademyFavorScheduleExcel struct {
	AcademyFavorScheduleExcelMap map[int64]*sro.AcademyFavorScheduleExcelTable
}

func (g *GameConfig) gppAcademyFavorScheduleExcelTable() {
	g.GetGPP().AcademyFavorScheduleExcel = &AcademyFavorScheduleExcel{
		AcademyFavorScheduleExcelMap: make(map[int64]*sro.AcademyFavorScheduleExcelTable),
	}
	for _, v := range g.GetExcel().GetAcademyFavorScheduleExcelTable() {
		g.GetGPP().AcademyFavorScheduleExcel.AcademyFavorScheduleExcelMap[v.Id] = v
	}

	logger.Info("处理MomoTalk剧情配置完成,剧情:%v个",
		len(g.GetGPP().AcademyFavorScheduleExcel.AcademyFavorScheduleExcelMap))
}

func GetAcademyFavorScheduleExcelTable(id int64) *sro.AcademyFavorScheduleExcelTable {
	return GC.GetGPP().AcademyFavorScheduleExcel.AcademyFavorScheduleExcelMap[id]
}
