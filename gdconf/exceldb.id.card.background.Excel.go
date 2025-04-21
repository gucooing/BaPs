package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadIdCardBackgroundExcel() {
	g.GetExcel().IdCardBackgroundExcel = make([]*sro.IdCardBackgroundExcel, 0)
	name := "IdCardBackgroundExcel.json"
	loadExcelJson(g.excelDbPath+name, &g.GetExcel().IdCardBackgroundExcel)
}

type IdCardBackground struct {
	IdCardBackgroundMap map[int64]*sro.IdCardBackgroundExcel
}

func (g *GameConfig) gppIdCardBackgroundExcel() {
	g.GetGPP().IdCardBackground = &IdCardBackground{
		IdCardBackgroundMap: make(map[int64]*sro.IdCardBackgroundExcel),
	}
	for _, v := range g.GetExcel().GetIdCardBackgroundExcel() {
		g.GetGPP().IdCardBackground.IdCardBackgroundMap[v.Id] = v
	}

	logger.Info("处理账号背景配置完成,数量:%v个",
		len(g.GetGPP().IdCardBackground.IdCardBackgroundMap))
}

func GetIdCardBackgroundExcelList() []*sro.IdCardBackgroundExcel {
	return GC.GetExcel().IdCardBackgroundExcel
}

func GetIdCardBackgroundExcel(id int64) *sro.IdCardBackgroundExcel {
	return GC.GetGPP().IdCardBackground.IdCardBackgroundMap[id]
}
