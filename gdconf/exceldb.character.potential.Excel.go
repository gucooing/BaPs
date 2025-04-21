package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadCharacterPotentialExcel() {
	g.GetExcel().CharacterPotentialExcel = make([]*sro.CharacterPotentialExcel, 0)
	name := "CharacterPotentialExcel.json"
	mx.LoadExcelJson(g.excelDbPath+name, &g.GetExcel().CharacterPotentialExcel)
}

type CharacterPotential struct {
	CharacterPotentialList          map[int64]*sro.CharacterPotentialExcel
	CharacterPotentialByCharacterId map[int64]map[string]*sro.CharacterPotentialExcel
}

func (g *GameConfig) gppCharacterPotentialExcel() {
	g.GetGPP().CharacterPotential = &CharacterPotential{
		CharacterPotentialList:          make(map[int64]*sro.CharacterPotentialExcel),
		CharacterPotentialByCharacterId: make(map[int64]map[string]*sro.CharacterPotentialExcel),
	}
	for _, v := range g.GetExcel().GetCharacterPotentialExcel() {
		g.GetGPP().CharacterPotential.CharacterPotentialList[v.PotentialStatGroupId] = v
		if g.GetGPP().CharacterPotential.CharacterPotentialByCharacterId[v.Id] == nil {
			g.GetGPP().CharacterPotential.CharacterPotentialByCharacterId[v.Id] = make(map[string]*sro.CharacterPotentialExcel)
		}
		g.GetGPP().CharacterPotential.CharacterPotentialByCharacterId[v.Id][v.PotentialStatBonusRateType] = v
	}

	logger.Info("处理角色能力解放配置完成,数量:%v个",
		len(g.GetGPP().CharacterPotential.CharacterPotentialByCharacterId))
}

func GetCharacterPotentialExcelType(characterId int64, rateType string) *sro.CharacterPotentialExcel {
	if GC.GetGPP().CharacterPotential.CharacterPotentialByCharacterId[characterId] == nil {
		return nil
	}
	return GC.GetGPP().CharacterPotential.CharacterPotentialByCharacterId[characterId][rateType]
}
