package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyLocationExcelTable() {
	g.GetExcel().AcademyLocationExcelTable = make([]*sro.AcademyLocationExcelTable, 0)
	name := "AcademyLocationExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().AcademyLocationExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetAcademyLocationExcelTable()))
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
