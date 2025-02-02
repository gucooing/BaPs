package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadRaidStageRewardExcelTable() {
	g.GetExcel().RaidStageRewardExcelTable = make([]*sro.RaidStageRewardExcelTable, 0)
	name := "RaidStageRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().RaidStageRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetRaidStageRewardExcelTable()))
}

type RaidStageRewardExcel struct {
	RaidStageRewardExcelMap map[int64][]*sro.RaidStageRewardExcelTable
}

func (g *GameConfig) gppRaidStageRewardExcelTable() {
	g.GetGPP().RaidStageRewardExcel = &RaidStageRewardExcel{
		RaidStageRewardExcelMap: make(map[int64][]*sro.RaidStageRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetRaidStageRewardExcelTable() {
		if g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[v.GroupId] == nil {
			g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[v.GroupId] = make([]*sro.RaidStageRewardExcelTable, 0)
		}
		g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[v.GroupId] =
			append(g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[v.GroupId], v)
	}

	logger.Info("处理总力战关卡通过奖励配置表完成,总力战关卡通过奖励配置:%v个",
		len(g.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap))
}

func GetRaidStageRewardExcelTable(gid int64) []*sro.RaidStageRewardExcelTable {
	return GC.GetGPP().RaidStageRewardExcel.RaidStageRewardExcelMap[gid]
}
