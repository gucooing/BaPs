package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadRaidRankingRewardExcelTable() {
	g.GetExcel().RaidRankingRewardExcelTable = make([]*sro.RaidRankingRewardExcelTable, 0)
	name := "RaidRankingRewardExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().RaidRankingRewardExcelTable)
}

type RaidRankingRewardExcel struct {
	RaidRankingRewardExcelMap map[int64][]*sro.RaidRankingRewardExcelTable
}

func (g *GameConfig) gppRaidRankingRewardExcelTable() {
	g.GetGPP().RaidRankingRewardExcel = &RaidRankingRewardExcel{
		RaidRankingRewardExcelMap: make(map[int64][]*sro.RaidRankingRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetRaidRankingRewardExcelTable() {
		if g.GetGPP().RaidRankingRewardExcel.RaidRankingRewardExcelMap[v.RankingRewardGroupId] == nil {
			g.GetGPP().RaidRankingRewardExcel.RaidRankingRewardExcelMap[v.RankingRewardGroupId] =
				make([]*sro.RaidRankingRewardExcelTable, 0)
		}
		g.GetGPP().RaidRankingRewardExcel.RaidRankingRewardExcelMap[v.RankingRewardGroupId] =
			append(g.GetGPP().RaidRankingRewardExcel.RaidRankingRewardExcelMap[v.RankingRewardGroupId], v)
	}

	logger.Info("处理总力战结算奖励配置表完成,奖励配置:%v个",
		len(g.GetGPP().RaidRankingRewardExcel.RaidRankingRewardExcelMap))
}

func GetRaidRankingRewardExcelTable(gid, ranking int64) *sro.RaidRankingRewardExcelTable {
	for _, conf := range GC.GetGPP().RaidRankingRewardExcel.RaidRankingRewardExcelMap[gid] {
		if conf.RankStart <= ranking && (conf.RankEnd >= ranking || conf.RankEnd == 0) {
			return conf
		}
		if ranking <= 0 && conf.RankEnd == 0 {
			return conf
		}
	}
	return nil
}
