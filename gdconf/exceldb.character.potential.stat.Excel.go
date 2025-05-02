package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadCharacterPotentialStatExcel() {
	g.GetExcel().CharacterPotentialStatExcel = make([]*sro.CharacterPotentialStatExcel, 0)
	name := "CharacterPotentialStatExcel.json"
	mx.LoadExcelJson(g.excelDbPath+name, &g.GetExcel().CharacterPotentialStatExcel)
}

type CharacterPotentialStat struct {
	CharacterPotentialList map[int64]map[int32]*sro.CharacterPotentialStatExcel
}

func (g *GameConfig) gppCharacterPotentialStatExcel() {
	g.GetGPP().CharacterPotentialStat = &CharacterPotentialStat{
		CharacterPotentialList: make(map[int64]map[int32]*sro.CharacterPotentialStatExcel),
	}
	for _, v := range g.GetExcel().GetCharacterPotentialStatExcel() {
		if g.GetGPP().CharacterPotentialStat.CharacterPotentialList[v.PotentialStatGroupId] == nil {
			g.GetGPP().CharacterPotentialStat.CharacterPotentialList[v.PotentialStatGroupId] = make(map[int32]*sro.CharacterPotentialStatExcel)
		}
		g.GetGPP().CharacterPotentialStat.CharacterPotentialList[v.PotentialStatGroupId][v.PotentialLevel] = v
	}

	logger.Info("处理角色能力解放升级配置完成,数量:%v个",
		len(g.GetGPP().CharacterPotentialStat.CharacterPotentialList))
}

func GetCharacterPotentialStatExcel(guid int64, level int32) *sro.CharacterPotentialStatExcel {
	if GC.GetGPP().CharacterPotentialStat.CharacterPotentialList[guid] == nil {
		return nil
	}
	return GC.GetGPP().CharacterPotentialStat.CharacterPotentialList[guid][level]
}
