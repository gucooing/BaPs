package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadDefaultEchelonExcelTable() {
	g.GetExcel().DefaultEchelonExcelTable = make([]*sro.DefaultEchelonExcelTable, 0)
	name := "DefaultEchelonExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().DefaultEchelonExcelTable)
}

func GetDefaultEchelonExcelList() []*sro.DefaultEchelonExcelTable {
	return GC.GetExcel().GetDefaultEchelonExcelTable()
}
