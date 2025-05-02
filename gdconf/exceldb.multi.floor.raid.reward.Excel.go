package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadMultiFloorRaidRewardExcel() {
	g.GetExcel().MultiFloorRaidRewardExcel = make([]*sro.MultiFloorRaidRewardExcel, 0)
	name := "MultiFloorRaidRewardExcel.json"
	mx.LoadExcelJson(g.excelDbPath+name, &g.GetExcel().MultiFloorRaidRewardExcel)
}

type MultiFloorRaidReward struct {
	MultiFloorRaidRewardMap map[int64][]*sro.MultiFloorRaidRewardExcel
}

func (g *GameConfig) gppMultiFloorRaidRewardExcel() {
	g.GetGPP().MultiFloorRaidReward = &MultiFloorRaidReward{
		MultiFloorRaidRewardMap: make(map[int64][]*sro.MultiFloorRaidRewardExcel),
	}
	for _, v := range g.GetExcel().GetMultiFloorRaidRewardExcel() {
		if g.GetGPP().MultiFloorRaidReward.MultiFloorRaidRewardMap[v.RewardGroupId] == nil {
			g.GetGPP().MultiFloorRaidReward.MultiFloorRaidRewardMap[v.RewardGroupId] =
				make([]*sro.MultiFloorRaidRewardExcel, 0)
		}
		g.GetGPP().MultiFloorRaidReward.MultiFloorRaidRewardMap[v.RewardGroupId] = append(
			g.GetGPP().MultiFloorRaidReward.MultiFloorRaidRewardMap[v.RewardGroupId], v)
	}

	logger.Info("处理制约解除决战关卡奖励完成,制约解除决战关卡奖励:%v个",
		len(g.GetGPP().MultiFloorRaidReward.MultiFloorRaidRewardMap))
}

func GetMultiFloorRaidRewardExcel(gid int64) []*sro.MultiFloorRaidRewardExcel {
	return GC.GetGPP().MultiFloorRaidReward.MultiFloorRaidRewardMap[gid]
}

func GetMultiFloorRaidRewardExcelBySeasonId(seasonId int64, difficulty int32) []*sro.MultiFloorRaidRewardExcel {
	conf := GetMultiFloorRaidStageExcelBySeason(seasonId, difficulty)
	if conf == nil {
		return nil
	}
	return GC.GetGPP().MultiFloorRaidReward.MultiFloorRaidRewardMap[conf.RewardGroupId]
}
