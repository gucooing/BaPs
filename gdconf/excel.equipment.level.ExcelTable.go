package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadEquipmentLevelExcelTable() {
	g.GetExcel().EquipmentLevelExcelTable = make([]*sro.EquipmentLevelExcelTable, 0)
	name := "EquipmentLevelExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().EquipmentLevelExcelTable)
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
