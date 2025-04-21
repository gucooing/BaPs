package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyTicketExcelTable() {
	g.GetExcel().AcademyTicketExcelTable = make([]*sro.AcademyTicketExcelTable, 0)
	name := "AcademyTicketExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().AcademyTicketExcelTable)
}

type AcademyTicketExcel struct {
	AcademyTicketExcelMap map[int64]*sro.AcademyTicketExcelTable
}

func (g *GameConfig) gppAcademyTicketExcelTable() {
	g.GetGPP().AcademyTicketExcel = &AcademyTicketExcel{
		AcademyTicketExcelMap: make(map[int64]*sro.AcademyTicketExcelTable),
	}
	for _, v := range g.GetExcel().GetAcademyTicketExcelTable() {
		g.GetGPP().AcademyTicketExcel.AcademyTicketExcelMap[v.ScheduleTicktetMax] = v
	}

	logger.Info("处理课程表最大票信息完成,数量:%v个",
		len(g.GetGPP().AcademyTicketExcel.AcademyTicketExcelMap))
}

func GetScheduleTicktetMax(level int64) int64 {
	for i := int64(3); ; i++ {
		conf := GC.GetGPP().AcademyTicketExcel.AcademyTicketExcelMap[i]
		if conf == nil {
			return i - 1
		}
		if conf.LocationRankSum > level {
			return i - 1
		}
	}
}
