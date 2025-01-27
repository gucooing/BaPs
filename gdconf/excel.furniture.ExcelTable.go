package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadFurnitureExcelTable() {
	g.GetExcel().FurnitureExcelTable = make([]*sro.FurnitureExcelTable, 0)
	name := "FurnitureExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().FurnitureExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetFurnitureExcelTable()))
}

type FurnitureExcel struct {
	FurnitureExcelTableMap map[int64]*sro.FurnitureExcelTable
}

func (g *GameConfig) gppFurnitureExcelTable() {
	g.GetGPP().FurnitureExcel = &FurnitureExcel{
		FurnitureExcelTableMap: make(map[int64]*sro.FurnitureExcelTable, 0),
	}
	for _, v := range g.GetExcel().GetFurnitureExcelTable() {
		g.GetGPP().FurnitureExcel.FurnitureExcelTableMap[v.Id] = v
	}
	logger.Info("处理家具配置表完成数量:%v个", len(g.GetGPP().FurnitureExcel.FurnitureExcelTableMap))
}

func GetFurnitureExcelTable(id int64) *sro.FurnitureExcelTable {
	return GC.GetGPP().FurnitureExcel.FurnitureExcelTableMap[id]
}

func GetFurnitureExcelTableMap() []*sro.FurnitureExcelTable {
	return GC.GetExcel().GetFurnitureExcelTable()
}
