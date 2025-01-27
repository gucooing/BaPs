package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCharacterPotentialExcel() {
	g.GetExcel().CharacterPotentialExcel = make([]*sro.CharacterPotentialExcel, 0)
	name := "CharacterPotentialExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CharacterPotentialExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().CharacterPotentialExcel))
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
