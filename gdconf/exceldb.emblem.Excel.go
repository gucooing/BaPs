package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEmblemExcel() {
	g.GetExcel().EmblemExcel = make([]*sro.EmblemExcel, 0)
	name := "EmblemExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EmblemExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().EmblemExcel))
}

type Emblem struct {
	EmblemExcelMap         map[int64]*sro.EmblemExcel
	EmblemExcelCategoryMap map[string][]*sro.EmblemExcel
}

func (g *GameConfig) gppEmblemExcel() {
	g.GetGPP().Emblem = &Emblem{
		EmblemExcelMap:         make(map[int64]*sro.EmblemExcel),
		EmblemExcelCategoryMap: make(map[string][]*sro.EmblemExcel),
	}
	for _, v := range g.GetExcel().GetEmblemExcel() {
		g.GetGPP().Emblem.EmblemExcelMap[v.Id] = v
		if g.GetGPP().Emblem.EmblemExcelCategoryMap[v.Category] == nil {
			g.GetGPP().Emblem.EmblemExcelCategoryMap[v.Category] = make([]*sro.EmblemExcel, 0)
		}
		g.GetGPP().Emblem.EmblemExcelCategoryMap[v.Category] = append(
			g.GetGPP().Emblem.EmblemExcelCategoryMap[v.Category], v)
	}

	logger.Info("处理称号配置完成,称号:%v个,获取类型:%v个", len(g.GetGPP().Emblem.EmblemExcelMap),
		len(g.GetGPP().Emblem.EmblemExcelCategoryMap))
}

func GetEmblemExcelList() []*sro.EmblemExcel {
	return GC.GetExcel().GetEmblemExcel()
}

func GetEmblemExcelCategoryList(category string) []*sro.EmblemExcel {
	return GC.GetGPP().Emblem.EmblemExcelCategoryMap[category]
}

func GetEmblemExcel(id int64) *sro.EmblemExcel {
	return GC.GetGPP().Emblem.EmblemExcelMap[id]
}
