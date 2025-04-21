package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadScenarioModeExcel() {
	g.GetExcel().ScenarioModeExcel = make([]*sro.ScenarioModeExcel, 0)
	name := "ScenarioModeExcel.json"
	loadExcelJson(g.excelDbPath+name, &g.GetExcel().ScenarioModeExcel)
}

type ScenarioMode struct {
	ScenarioModeMap map[int64]*sro.ScenarioModeExcel
}

func (g *GameConfig) gppScenarioModeExcel() {
	g.GetGPP().ScenarioMode = &ScenarioMode{
		ScenarioModeMap: make(map[int64]*sro.ScenarioModeExcel),
	}
	for _, v := range g.GetExcel().GetScenarioModeExcel() {
		g.GetGPP().ScenarioMode.ScenarioModeMap[v.ModeId] = v
	}

	logger.Info("处理剧情配置完成,剧情:%v个",
		len(g.GetGPP().ScenarioMode.ScenarioModeMap))
}

func GetScenarioModeExcel(id int64) *sro.ScenarioModeExcel {
	return GC.GetGPP().ScenarioMode.ScenarioModeMap[id]
}
