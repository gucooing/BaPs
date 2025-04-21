package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadMultiFloorRaidStageExcel() {
	g.GetExcel().MultiFloorRaidStageExcel = make([]*sro.MultiFloorRaidStageExcel, 0)
	name := "MultiFloorRaidStageExcel.json"
	mx.LoadExcelJson(g.excelDbPath+name, &g.GetExcel().MultiFloorRaidStageExcel)
}

type MultiFloorRaidStage struct {
	MultiFloorRaidStageMap map[string]map[int32]*sro.MultiFloorRaidStageExcel
}

func (g *GameConfig) gppMultiFloorRaidStageExcel() {
	g.GetGPP().MultiFloorRaidStage = &MultiFloorRaidStage{
		MultiFloorRaidStageMap: make(map[string]map[int32]*sro.MultiFloorRaidStageExcel),
	}
	for _, v := range g.GetExcel().GetMultiFloorRaidStageExcel() {
		if g.GetGPP().MultiFloorRaidStage.MultiFloorRaidStageMap[v.BossGroupId] == nil {
			g.GetGPP().MultiFloorRaidStage.MultiFloorRaidStageMap[v.BossGroupId] =
				make(map[int32]*sro.MultiFloorRaidStageExcel)
		}
		g.GetGPP().MultiFloorRaidStage.MultiFloorRaidStageMap[v.BossGroupId][v.Difficulty] = v
	}

	logger.Info("处理制约解除决战关卡完成,制约解除决战关卡:%v个",
		len(g.GetGPP().MultiFloorRaidStage.MultiFloorRaidStageMap))
}

func GetMultiFloorRaidStageExcel(bgid string, difficulty int32) *sro.MultiFloorRaidStageExcel {
	conf := GC.GetGPP().MultiFloorRaidStage.MultiFloorRaidStageMap[bgid]
	if conf == nil {
		return nil
	}
	return conf[difficulty]
}

func GetMultiFloorRaidStageExcelBySeason(seasonId int64, difficulty int32) *sro.MultiFloorRaidStageExcel {
	conf := GetMultiFloorRaidSeasonManageExcel(seasonId)
	if conf == nil {
		return nil
	}
	return GetMultiFloorRaidStageExcel(conf.OpenRaidBossGroupId, difficulty)
}
