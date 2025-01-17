package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyFavorScheduleExcelTable() {
	g.GetExcel().AcademyFavorScheduleExcelTable = make([]*sro.AcademyFavorScheduleExcelTable, 0)
	name := "AcademyFavorScheduleExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().AcademyFavorScheduleExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().AcademyFavorScheduleExcelTable))
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
