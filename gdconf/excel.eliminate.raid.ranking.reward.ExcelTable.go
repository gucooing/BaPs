package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadEliminateRaidRankingRewardExcelTable() {
	g.GetExcel().EliminateRaidRankingRewardExcelTable = make([]*sro.EliminateRaidRankingRewardExcelTable, 0)
	name := "EliminateRaidRankingRewardExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().EliminateRaidRankingRewardExcelTable)
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
