package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadIdCardBackgroundExcel() {
	g.GetExcel().IdCardBackgroundExcel = make([]*sro.IdCardBackgroundExcel, 0)
	name := "IdCardBackgroundExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().IdCardBackgroundExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().IdCardBackgroundExcel))
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
