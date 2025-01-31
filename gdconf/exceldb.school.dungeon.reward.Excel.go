package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadSchoolDungeonRewardExcel() {
	g.GetExcel().SchoolDungeonRewardExcel = make([]*sro.SchoolDungeonRewardExcel, 0)
	name := "SchoolDungeonRewardExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().SchoolDungeonRewardExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetSchoolDungeonRewardExcel()))
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
