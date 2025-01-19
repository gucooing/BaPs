package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadScenarioModeRewardExcel() {
	g.GetExcel().ScenarioModeRewardExcel = make([]*sro.ScenarioModeRewardExcel, 0)
	name := "ScenarioModeRewardExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().ScenarioModeRewardExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().ScenarioModeRewardExcel))
}

type ScenarioModeReward struct {
	ScenarioModeRewardMap map[int64][]*sro.ScenarioModeRewardExcel
}

func (g *GameConfig) gppScenarioModeRewardExcel() {
	g.GetGPP().ScenarioModeReward = &ScenarioModeReward{
		ScenarioModeRewardMap: make(map[int64][]*sro.ScenarioModeRewardExcel),
	}
	for _, v := range g.GetExcel().GetScenarioModeRewardExcel() {
		if g.GetGPP().ScenarioModeReward.ScenarioModeRewardMap[v.ScenarioModeRewardId] == nil {
			g.GetGPP().ScenarioModeReward.ScenarioModeRewardMap[v.ScenarioModeRewardId] =
				make([]*sro.ScenarioModeRewardExcel, 0)
		}
		g.GetGPP().ScenarioModeReward.ScenarioModeRewardMap[v.ScenarioModeRewardId] = append(
			g.GetGPP().ScenarioModeReward.ScenarioModeRewardMap[v.ScenarioModeRewardId], v)
	}

	logger.Info("处理剧情奖励配置完成,剧情奖励:%v个",
		len(g.GetGPP().ScenarioModeReward.ScenarioModeRewardMap))
}

func GetScenarioModeRewardExcel(id int64) []*sro.ScenarioModeRewardExcel {
	return GC.GetGPP().ScenarioModeReward.ScenarioModeRewardMap[id]
}
