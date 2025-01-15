package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadShopExcelTable() {
	g.GetExcel().ShopExcelTable = make([]*sro.ShopExcelTable, 0)
	name := "ShopExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().ShopExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetShopExcelTable()))
}

type ShopExcel struct {
	ShopExcelTableMap map[int64]*sro.ShopExcelTable
	ShopExcelTypeMap  map[string][]*sro.ShopExcelTable
}

func (g *GameConfig) gppShopExcelTable() {
	g.GetGPP().ShopExcel = &ShopExcel{
		ShopExcelTableMap: make(map[int64]*sro.ShopExcelTable),
		ShopExcelTypeMap:  make(map[string][]*sro.ShopExcelTable),
	}
	for _, v := range g.GetExcel().GetShopExcelTable() {
		g.GetGPP().ShopExcel.ShopExcelTableMap[v.Id] = v
		if g.GetGPP().ShopExcel.ShopExcelTypeMap[v.CategoryType] == nil {
			g.GetGPP().ShopExcel.ShopExcelTypeMap[v.CategoryType] = make([]*sro.ShopExcelTable, 0)
		}
		g.GetGPP().ShopExcel.ShopExcelTypeMap[v.CategoryType] = append(
			g.GetGPP().ShopExcel.ShopExcelTypeMap[v.CategoryType], v)
	}
	logger.Info("处理商品配置完成,商店类型:%v个", len(g.GetGPP().ShopExcel.ShopExcelTypeMap))
}

func GetShopExcelType(categoryType string) []*sro.ShopExcelTable {
	if GC.GetGPP().ShopExcel.ShopExcelTypeMap == nil {
		return nil
	}
	return GC.GetGPP().ShopExcel.ShopExcelTypeMap[categoryType]
}
