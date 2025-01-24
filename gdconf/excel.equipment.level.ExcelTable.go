package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEquipmentLevelExcelTable() {
	g.GetExcel().EquipmentLevelExcelTable = make([]*sro.EquipmentLevelExcelTable, 0)
	name := "EquipmentLevelExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EquipmentLevelExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetEquipmentLevelExcelTable()))
}

type EquipmentLevelExcel struct {
	EquipmentLevelExcelTableMap map[int32]*sro.EquipmentLevelExcelTable
}

func (g *GameConfig) gppEquipmentLevelExcelTable() {
	g.GetGPP().EquipmentLevelExcel = &EquipmentLevelExcel{
		EquipmentLevelExcelTableMap: make(map[int32]*sro.EquipmentLevelExcelTable, 0),
	}
	for _, v := range g.GetExcel().GetEquipmentLevelExcelTable() {
		g.GetGPP().EquipmentLevelExcel.EquipmentLevelExcelTableMap[v.Level] = v
	}
	logger.Info("装备等级配置表完成数量:%v个", len(g.GetGPP().EquipmentLevelExcel.EquipmentLevelExcelTableMap))
}

func GetEquipmentLevelExcelTable(level int32) *sro.EquipmentLevelExcelTable {
	return GC.GetGPP().EquipmentLevelExcel.EquipmentLevelExcelTableMap[level]
}
