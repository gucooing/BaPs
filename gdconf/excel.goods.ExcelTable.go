package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadGoodsExcelTable() {
	g.GetExcel().GoodsExcelTable = make([]*sro.GoodsExcelTable, 0)
	name := "GoodsExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().GoodsExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GoodsExcelTable))
}

type GoodsExcel struct {
	GoodsExcelMap map[int64]*sro.GoodsExcelTable
}

func (g *GameConfig) gppGoodsExcelTable() {
	g.GetGPP().GoodsExcel = &GoodsExcel{
		GoodsExcelMap: make(map[int64]*sro.GoodsExcelTable),
	}
	for _, v := range g.GetExcel().GetGoodsExcelTable() {
		g.GetGPP().GoodsExcel.GoodsExcelMap[v.Id] = v
	}

	logger.Info("处理商品配置完成,成就:%v个",
		len(g.GetGPP().GoodsExcel.GoodsExcelMap))
}

func GetGoodsExcelTable(id int64) *sro.GoodsExcelTable {
	return GC.GetGPP().GoodsExcel.GoodsExcelMap[id]
}
