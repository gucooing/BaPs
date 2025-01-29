package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyZoneExcelTable() {
	g.GetExcel().AcademyZoneExcelTable = make([]*sro.AcademyZoneExcelTable, 0)
	name := "AcademyZoneExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().AcademyZoneExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetAcademyZoneExcelTable()))
}

type AcademyZoneExcel struct {
	AcademyZoneExcelMap map[int64]*sro.AcademyZoneExcelTable
}

func (g *GameConfig) gppAcademyZoneExcelTable() {
	g.GetGPP().AcademyZoneExcel = &AcademyZoneExcel{
		AcademyZoneExcelMap: make(map[int64]*sro.AcademyZoneExcelTable),
	}
	for _, v := range g.GetExcel().GetAcademyZoneExcelTable() {
		g.GetGPP().AcademyZoneExcel.AcademyZoneExcelMap[v.Id] = v
	}

	logger.Info("处理课程表教室信息完成,数量:%v个",
		len(g.GetGPP().AcademyZoneExcel.AcademyZoneExcelMap))
}

func GetAcademyZoneExcelTableList() []*sro.AcademyZoneExcelTable {
	return GC.GetExcel().GetAcademyZoneExcelTable()
}

func GetAcademyZoneExcelTable(zoneId int64) *sro.AcademyZoneExcelTable {
	return GC.GetGPP().AcademyZoneExcel.AcademyZoneExcelMap[zoneId]
}
