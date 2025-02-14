package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEliminateRaidStageExcelTable() {
	g.GetExcel().EliminateRaidStageExcelTable = make([]*sro.EliminateRaidStageExcelTable, 0)
	name := "EliminateRaidStageExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EliminateRaidStageExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetEliminateRaidStageExcelTable()))
}

type EliminateRaidStageExcel struct {
	EliminateRaidStageExcelMap map[int64]*sro.EliminateRaidStageExcelTable
}

func (g *GameConfig) gppEliminateRaidStageExcelTable() {
	g.GetGPP().EliminateRaidStageExcel = &EliminateRaidStageExcel{
		EliminateRaidStageExcelMap: make(map[int64]*sro.EliminateRaidStageExcelTable),
	}

	for _, v := range g.GetExcel().GetEliminateRaidStageExcelTable() {
		g.GetGPP().EliminateRaidStageExcel.EliminateRaidStageExcelMap[v.Id] = v
	}

	logger.Info("处理大决战关卡配置表完成,关卡配置:%v个",
		len(g.GetGPP().EliminateRaidStageExcel.EliminateRaidStageExcelMap))
}

func GetEliminateRaidStageExcelTable(id int64) *sro.EliminateRaidStageExcelTable {
	return GC.GetGPP().EliminateRaidStageExcel.EliminateRaidStageExcelMap[id]
}
