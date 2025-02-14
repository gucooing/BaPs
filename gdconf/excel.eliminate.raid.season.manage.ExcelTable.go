package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEliminateRaidSeasonManageExcelTable() {
	g.GetExcel().EliminateRaidSeasonManageExcelTable = make([]*sro.EliminateRaidSeasonManageExcelTable, 0)
	name := "EliminateRaidSeasonManageExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EliminateRaidSeasonManageExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetEliminateRaidSeasonManageExcelTable()))
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
