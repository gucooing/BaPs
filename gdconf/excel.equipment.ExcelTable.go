package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEquipmentExcelTable() {
	g.GetExcel().EquipmentExcelTable = make([]*sro.EquipmentExcelTable, 0)
	name := "EquipmentExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EquipmentExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetEquipmentExcelTable()))
}

type EquipmentExcel struct {
	EquipmentExcelTableMap map[int64]*sro.EquipmentExcelTable
}

func (g *GameConfig) gppEquipmentExcelTable() {
	g.GetGPP().EquipmentExcel = &EquipmentExcel{
		EquipmentExcelTableMap: make(map[int64]*sro.EquipmentExcelTable, 0),
	}
	for _, v := range g.GetExcel().GetEquipmentExcelTable() {
		g.GetGPP().EquipmentExcel.EquipmentExcelTableMap[v.Id] = v
	}
	logger.Info("装备配置表完成数量:%v个", len(g.GetGPP().EquipmentExcel.EquipmentExcelTableMap))
}

func GetEquipmentExcelTableMap() []*sro.EquipmentExcelTable {
	return GC.GetExcel().GetEquipmentExcelTable()
}

func GetEquipmentExcelTable(id int64) *sro.EquipmentExcelTable {
	return GC.GetGPP().EquipmentExcel.EquipmentExcelTableMap[id]
}
