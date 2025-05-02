package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadDefaultCharacterExcelTable() {
	g.GetExcel().DefaultCharacterExcelTable = make([]*sro.DefaultCharacterExcelTable, 0)
	name := "DefaultCharacterExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().DefaultCharacterExcelTable)
}

func GetDefaultCharacterExcelTable() []*sro.DefaultCharacterExcelTable {
	if e := GC.GetExcel(); e == nil {
		return nil
	} else {
		return e.GetDefaultCharacterExcelTable()
	}
}
