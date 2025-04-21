package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadCharacterWeaponExcelTable() {
	g.GetExcel().CharacterWeaponExcelTable = make([]*sro.CharacterWeaponExcelTable, 0)
	name := "CharacterWeaponExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().CharacterWeaponExcelTable)
}

type CharacterWeaponExcel struct {
	CharacterWeaponExcelMap map[int64]*sro.CharacterWeaponExcelTable
}

func (g *GameConfig) gppCharacterWeaponExcelTable() {
	g.GetGPP().CharacterWeaponExcel = &CharacterWeaponExcel{
		CharacterWeaponExcelMap: make(map[int64]*sro.CharacterWeaponExcelTable),
	}

	for _, v := range g.GetExcel().GetCharacterWeaponExcelTable() {
		g.GetGPP().CharacterWeaponExcel.CharacterWeaponExcelMap[v.Id] = v
	}

	logger.Info("角色武器配置完成,角色武器:%v个",
		len(g.GetGPP().CharacterWeaponExcel.CharacterWeaponExcelMap))
}

func GetCharacterWeaponExcelTable(characterId int64) *sro.CharacterWeaponExcelTable {
	return GC.GetGPP().CharacterWeaponExcel.CharacterWeaponExcelMap[characterId]
}
