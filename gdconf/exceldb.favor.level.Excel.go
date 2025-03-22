package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadFavorLevelExcel() {
	g.GetExcel().FavorLevelExcel = make([]*sro.FavorLevelExcel, 0)
	name := "FavorLevelExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().FavorLevelExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().FavorLevelExcel))
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
