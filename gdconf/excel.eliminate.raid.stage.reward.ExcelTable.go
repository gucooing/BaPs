package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEliminateRaidStageRewardExcelTable() {
	g.GetExcel().EliminateRaidStageRewardExcelTable = make([]*sro.EliminateRaidStageRewardExcelTable, 0)
	name := "EliminateRaidStageRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EliminateRaidStageRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetEliminateRaidStageRewardExcelTable()))
}

type EliminateRaidStageRewardExcel struct {
	EliminateRaidStageRewardExcelMap map[int64][]*sro.EliminateRaidStageRewardExcelTable
}

func (g *GameConfig) gppEliminateRaidStageRewardExcelTable() {
	g.GetGPP().EliminateRaidStageRewardExcel = &EliminateRaidStageRewardExcel{
		EliminateRaidStageRewardExcelMap: make(map[int64][]*sro.EliminateRaidStageRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetEliminateRaidStageRewardExcelTable() {
		if g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[v.GroupId] == nil {
			g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[v.GroupId] = make([]*sro.EliminateRaidStageRewardExcelTable, 0)
		}
		g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[v.GroupId] =
			append(g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[v.GroupId], v)
	}

	logger.Info("处理大决战关卡通过奖励配置表完成,总力战关卡通过奖励配置:%v个",
		len(g.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap))
}

func GetEliminateRaidStageRewardExcelTable(gid int64) []*sro.EliminateRaidStageRewardExcelTable {
	return GC.GetGPP().EliminateRaidStageRewardExcel.EliminateRaidStageRewardExcelMap[gid]
}
