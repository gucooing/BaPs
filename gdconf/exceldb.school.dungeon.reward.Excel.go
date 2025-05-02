package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadSchoolDungeonRewardExcel() {
	g.GetExcel().SchoolDungeonRewardExcel = make([]*sro.SchoolDungeonRewardExcel, 0)
	name := "SchoolDungeonRewardExcel.json"
	mx.LoadExcelJson(g.excelDbPath+name, &g.GetExcel().SchoolDungeonRewardExcel)
}

type SchoolDungeonReward struct {
	SchoolDungeonRewardMap map[int64][]*sro.SchoolDungeonRewardExcel
}

func (g *GameConfig) gppSchoolDungeonRewardExcel() {
	g.GetGPP().SchoolDungeonReward = &SchoolDungeonReward{
		SchoolDungeonRewardMap: make(map[int64][]*sro.SchoolDungeonRewardExcel, 0),
	}

	for _, v := range g.GetExcel().GetSchoolDungeonRewardExcel() {
		if g.GetGPP().SchoolDungeonReward.SchoolDungeonRewardMap[v.GroupId] == nil {
			g.GetGPP().SchoolDungeonReward.SchoolDungeonRewardMap[v.GroupId] = make([]*sro.SchoolDungeonRewardExcel, 0)
		}
		g.GetGPP().SchoolDungeonReward.SchoolDungeonRewardMap[v.GroupId] =
			append(g.GetGPP().SchoolDungeonReward.SchoolDungeonRewardMap[v.GroupId], v)
	}

	logger.Info("处理学院交流会奖励配置完成,数量:%v个", len(g.GetGPP().SchoolDungeonReward.SchoolDungeonRewardMap))
}

func GetSchoolDungeonRewardExcelList(gId int64) []*sro.SchoolDungeonRewardExcel {
	return GC.GetGPP().SchoolDungeonReward.SchoolDungeonRewardMap[gId]
}
