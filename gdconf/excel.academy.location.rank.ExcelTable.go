package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyLocationRankExcelTable() {
	g.GetExcel().AcademyLocationRankExcelTable = make([]*sro.AcademyLocationRankExcelTable, 0)
	name := "AcademyLocationRankExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().AcademyLocationRankExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetAcademyLocationRankExcelTable()))
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
