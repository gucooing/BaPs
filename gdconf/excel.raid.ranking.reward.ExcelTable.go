package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadRaidRankingRewardExcelTable() {
	g.GetExcel().RaidRankingRewardExcelTable = make([]*sro.RaidRankingRewardExcelTable, 0)
	name := "RaidRankingRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().RaidRankingRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetRaidRankingRewardExcelTable()))
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
