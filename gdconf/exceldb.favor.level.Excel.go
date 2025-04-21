package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadFavorLevelExcel() {
	g.GetExcel().FavorLevelExcel = make([]*sro.FavorLevelExcel, 0)
	name := "FavorLevelExcel.json"
	loadExcelJson(g.excelDbPath+name, &g.GetExcel().FavorLevelExcel)
}

type FavorLevel struct {
	FavorLevelExcelMap map[int32]*sro.FavorLevelExcel
}

func (g *GameConfig) gppFavorLevelExcel() {
	g.GetGPP().FavorLevel = &FavorLevel{
		FavorLevelExcelMap: make(map[int32]*sro.FavorLevelExcel),
	}
	for _, v := range g.GetExcel().GetFavorLevelExcel() {
		g.GetGPP().FavorLevel.FavorLevelExcelMap[v.Level] = v
	}

	logger.Info("好感系统等级经验配置完成,数量:%v个", len(g.GetGPP().FavorLevel.FavorLevelExcelMap))
}

func GetFavorLevelExcel(level int32) *sro.FavorLevelExcel {
	return GC.GetGPP().FavorLevel.FavorLevelExcelMap[level]
}

func GetMaxFavorLevel() int32 {
	return GC.GetExcel().FavorLevelExcel[len(GC.GetExcel().FavorLevelExcel)-1].GetLevel()
}
