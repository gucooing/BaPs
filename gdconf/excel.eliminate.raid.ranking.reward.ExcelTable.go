package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEliminateRaidRankingRewardExcelTable() {
	g.GetExcel().EliminateRaidRankingRewardExcelTable = make([]*sro.EliminateRaidRankingRewardExcelTable, 0)
	name := "EliminateRaidRankingRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EliminateRaidRankingRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetEliminateRaidRankingRewardExcelTable()))
}

type EliminateRaidRankingRewardExcel struct {
	EliminateRaidRankingRewardExcelMap map[int64][]*sro.EliminateRaidRankingRewardExcelTable
}

func (g *GameConfig) gppEliminateRaidRankingRewardExcelTable() {
	g.GetGPP().EliminateRaidRankingRewardExcel = &EliminateRaidRankingRewardExcel{
		EliminateRaidRankingRewardExcelMap: make(map[int64][]*sro.EliminateRaidRankingRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetEliminateRaidRankingRewardExcelTable() {
		if g.GetGPP().EliminateRaidRankingRewardExcel.EliminateRaidRankingRewardExcelMap[v.RankingRewardGroupId] == nil {
			g.GetGPP().EliminateRaidRankingRewardExcel.EliminateRaidRankingRewardExcelMap[v.RankingRewardGroupId] =
				make([]*sro.EliminateRaidRankingRewardExcelTable, 0)
		}
		g.GetGPP().EliminateRaidRankingRewardExcel.EliminateRaidRankingRewardExcelMap[v.RankingRewardGroupId] =
			append(g.GetGPP().EliminateRaidRankingRewardExcel.EliminateRaidRankingRewardExcelMap[v.RankingRewardGroupId], v)
	}

	logger.Info("处理大决战结算奖励配置表完成,奖励配置:%v个",
		len(g.GetGPP().EliminateRaidRankingRewardExcel.EliminateRaidRankingRewardExcelMap))
}

func GetEliminateRaidRankingRewardExcelTable(gid, ranking int64) *sro.EliminateRaidRankingRewardExcelTable {
	for _, conf := range GC.GetGPP().EliminateRaidRankingRewardExcel.EliminateRaidRankingRewardExcelMap[gid] {
		if conf.RankStart <= ranking && (conf.RankEnd >= ranking || conf.RankEnd == 0) {
			return conf
		}
		if ranking <= 0 && conf.RankEnd == 0 {
			return conf
		}
	}
	return nil
}
