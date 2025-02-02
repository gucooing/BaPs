package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadRaidStageExcelTable() {
	g.GetExcel().RaidStageExcelTable = make([]*sro.RaidStageExcelTable, 0)
	name := "RaidStageExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().RaidStageExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetRaidStageExcelTable()))
}

type RaidStageExcel struct {
	RaidStageExcelMap map[int64]*sro.RaidStageExcelTable
}

func (g *GameConfig) gppRaidStageExcelTable() {
	g.GetGPP().RaidStageExcel = &RaidStageExcel{
		RaidStageExcelMap: make(map[int64]*sro.RaidStageExcelTable),
	}

	for _, v := range g.GetExcel().GetRaidStageExcelTable() {
		g.GetGPP().RaidStageExcel.RaidStageExcelMap[v.Id] = v
	}

	logger.Info("处理总力战关卡配置表完成,关卡配置:%v个",
		len(g.GetGPP().RaidStageExcel.RaidStageExcelMap))
}

func GetRaidStageExcelTable(id int64) *sro.RaidStageExcelTable {
	return GC.GetGPP().RaidStageExcel.RaidStageExcelMap[id]
}
