package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyMessangerExcelTable() {
	g.GetExcel().AcademyMessangerExcelTable = make([]*sro.AcademyMessangerExcelTable, 0)
	nameList := []string{"AcademyMessanger1ExcelTable.json", "AcademyMessanger2ExcelTable.json",
		"AcademyMessanger3ExcelTable.json", "AcademyMessangerExcelTable.json"}
	for _, name := range nameList {
		list := make([]*sro.AcademyMessangerExcelTable, 0)
		file, err := os.ReadFile(g.excelPath + name)
		if err != nil {
			logger.Error("文件:%s 读取失败,err:%s", name, err)
			return
		}
		if err := json.Unmarshal(file, &list); err != nil {
			logger.Error("文件:%s 解析失败,err:%s", name, err)
			return
		}
		g.GetExcel().AcademyMessangerExcelTable = append(g.GetExcel().AcademyMessangerExcelTable, list...)
		logger.Info("文件:%s 读取成功,解析数量:%v", name, len(list))
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
