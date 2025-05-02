package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadEquipmentStatExcelTable() {
	g.GetExcel().EquipmentStatExcelTable = make([]*sro.EquipmentStatExcelTable, 0)
	name := "EquipmentStatExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().EquipmentStatExcelTable)
}

type EquipmentStatExcel struct {
	EquipmentStatExcelTableMap map[int64]*sro.EquipmentStatExcelTable
}

func (g *GameConfig) gppEquipmentStatExcelTable() {
	g.GetGPP().EquipmentStatExcel = &EquipmentStatExcel{
		EquipmentStatExcelTableMap: make(map[int64]*sro.EquipmentStatExcelTable, 0),
	}
	for _, v := range g.GetExcel().GetEquipmentStatExcelTable() {
		g.GetGPP().EquipmentStatExcel.EquipmentStatExcelTableMap[v.EquipmentId] = v
	}
	logger.Info("装备详情表完成数量:%v个", len(g.GetGPP().EquipmentStatExcel.EquipmentStatExcelTableMap))
}

func GetEquipmentStatExcelTable(eId int64) *sro.EquipmentStatExcelTable {
	return GC.GetGPP().EquipmentStatExcel.EquipmentStatExcelTableMap[eId]
}
