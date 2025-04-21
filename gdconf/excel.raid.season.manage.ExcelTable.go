package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadRaidSeasonManageExcelTable() {
	g.GetExcel().RaidSeasonManageExcelTable = make([]*sro.RaidSeasonManageExcelTable, 0)
	name := "RaidSeasonManageExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().RaidSeasonManageExcelTable)
}

type RaidSeasonManageExcel struct {
	RaidSeasonManageExcelMap map[int64]*sro.RaidSeasonManageExcelTable
}

func (g *GameConfig) gppRaidSeasonManageExcelTable() {
	g.GetGPP().RaidSeasonManageExcel = &RaidSeasonManageExcel{
		RaidSeasonManageExcelMap: make(map[int64]*sro.RaidSeasonManageExcelTable),
	}

	for _, v := range g.GetExcel().GetRaidSeasonManageExcelTable() {
		g.GetGPP().RaidSeasonManageExcel.RaidSeasonManageExcelMap[v.SeasonId] = v
	}

	logger.Info("处理总力战赛季配置表完成,赛季配置:%v个",
		len(g.GetGPP().RaidSeasonManageExcel.RaidSeasonManageExcelMap))
}

func GetRaidSeasonManageExcelTable(seasonId int64) *sro.RaidSeasonManageExcelTable {
	return GC.GetGPP().RaidSeasonManageExcel.RaidSeasonManageExcelMap[seasonId]
}

func GetRaidTier(seasonId, ranking int64) int32 {
	conf := GetRaidSeasonManageExcelTable(seasonId)
	if conf == nil || ranking < 0 {
		return 1
	}
	return GetRaidRankingRewardExcelTable(conf.RankingRewardGroupId, ranking).GetTier()
}

func GetRaidRankingRewardExcelTableBySeasonId(seasonId, ranking int64) *sro.RaidRankingRewardExcelTable {
	conf := GetRaidSeasonManageExcelTable(seasonId)
	if conf == nil {
		return nil
	}
	return GetRaidRankingRewardExcelTable(conf.RankingRewardGroupId, ranking)
}
