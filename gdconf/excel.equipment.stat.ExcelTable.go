package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadEquipmentStatExcelTable() {
	g.GetExcel().EquipmentStatExcelTable = make([]*sro.EquipmentStatExcelTable, 0)
	name := "EquipmentStatExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().EquipmentStatExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetEquipmentStatExcelTable()))
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
