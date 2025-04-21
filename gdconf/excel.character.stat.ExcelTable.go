package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadCharacterStatExcelTable() {
	g.GetExcel().CharacterStatExcelTable = make([]*sro.CharacterStatExcelTable, 0)
	name := "CharacterStatExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().CharacterStatExcelTable)
}

type CharacterStatExcel struct {
	CharacterStatExcelMap map[int64]*sro.CharacterStatExcelTable
}

func (g *GameConfig) gppCharacterStatExcelTable() {
	g.GetGPP().CharacterStatExcel = &CharacterStatExcel{
		CharacterStatExcelMap: make(map[int64]*sro.CharacterStatExcelTable),
	}

	for _, v := range g.GetExcel().GetCharacterStatExcelTable() {
		g.GetGPP().CharacterStatExcel.CharacterStatExcelMap[v.CharacterId] = v
	}

	logger.Info("处理实体属性配置表完成,实体属性配置:%v个",
		len(g.GetGPP().CharacterStatExcel.CharacterStatExcelMap))
}

func GetCharacterStatExcelTable(characterId int64) *sro.CharacterStatExcelTable {
	return GC.GetGPP().CharacterStatExcel.CharacterStatExcelMap[characterId]
}
