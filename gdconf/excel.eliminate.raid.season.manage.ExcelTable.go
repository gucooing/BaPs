package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func (g *GameConfig) loadEliminateRaidSeasonManageExcelTable() {
	g.GetExcel().EliminateRaidSeasonManageExcelTable = make([]*sro.EliminateRaidSeasonManageExcelTable, 0)
	name := "EliminateRaidSeasonManageExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().EliminateRaidSeasonManageExcelTable)
}

type EliminateRaidSeasonManageExcel struct {
	EliminateRaidSeasonManageExcelMap map[int64]*sro.EliminateRaidSeasonManageExcelTable
}

func (g *GameConfig) gppEliminateRaidSeasonManageExcelTable() {
	g.GetGPP().EliminateRaidSeasonManageExcel = &EliminateRaidSeasonManageExcel{
		EliminateRaidSeasonManageExcelMap: make(map[int64]*sro.EliminateRaidSeasonManageExcelTable),
	}

	for _, v := range g.GetExcel().GetEliminateRaidSeasonManageExcelTable() {
		g.GetGPP().EliminateRaidSeasonManageExcel.EliminateRaidSeasonManageExcelMap[v.SeasonId] = v
	}

	logger.Info("处理大决战赛季配置表完成,赛季配置:%v个",
		len(g.GetGPP().EliminateRaidSeasonManageExcel.EliminateRaidSeasonManageExcelMap))
}

func GetEliminateRaidSeasonManageExcelTable(seasonId int64) *sro.EliminateRaidSeasonManageExcelTable {
	return GC.GetGPP().EliminateRaidSeasonManageExcel.EliminateRaidSeasonManageExcelMap[seasonId]
}

func GetEliminateRaidTier(seasonId, ranking int64) int32 {
	conf := GetEliminateRaidSeasonManageExcelTable(seasonId)
	if conf == nil || ranking < 0 {
		return 1
	}
	return GetEliminateRaidRankingRewardExcelTable(conf.RankingRewardGroupId, ranking).GetTier()
}

func GetEliminateRaidRankingRewardExcelTableBySeasonId(seasonId, ranking int64) *sro.EliminateRaidRankingRewardExcelTable {
	conf := GetEliminateRaidSeasonManageExcelTable(seasonId)
	if conf == nil {
		return nil
	}
	return GetEliminateRaidRankingRewardExcelTable(conf.RankingRewardGroupId, ranking)
}

func GetEliminateRaidEchelonType(seasonId int64, raidBossGroup string) proto.EchelonType {
	conf := GetEliminateRaidSeasonManageExcelTable(seasonId)
	if conf == nil {
		return proto.EchelonType_None
	}
	switch raidBossGroup {
	case conf.OpenRaidBossGroup01:
		return proto.EchelonType_EliminateRaid01
	case conf.OpenRaidBossGroup02:
		return proto.EchelonType_EliminateRaid02
	case conf.OpenRaidBossGroup03:
		return proto.EchelonType_EliminateRaid03
	default:
		return proto.EchelonType_None
	}
}
