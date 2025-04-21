package gdconf

import (
	"math/rand"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadCharacterExcelTable() {
	g.GetExcel().CharacterExcelTable = make([]*sro.CharacterExcelTable, 0)
	name := "CharacterExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().CharacterExcelTable)
}

type CharacterExcel struct {
	CharacterExcelTableMap map[int64]*sro.CharacterExcelTable
	CharacterReleaseList   []*sro.CharacterExcelTable
	CharacterMap           map[int64]*sro.CharacterExcelTable // 全部角色索引
	// CharacterRMap          []*sro.CharacterExcelTable
	// CharacterSRMap         []*sro.CharacterExcelTable
	// CharacterSSRMap        []*sro.CharacterExcelTable
}

/*
R 1
SR 2
SSR 3
*/

func (g *GameConfig) gppCharacterExcelTable() {
	info := &CharacterExcel{
		CharacterExcelTableMap: make(map[int64]*sro.CharacterExcelTable),
		CharacterReleaseList:   make([]*sro.CharacterExcelTable, 0),
		CharacterMap:           make(map[int64]*sro.CharacterExcelTable),
		// CharacterRMap:          make([]*sro.CharacterExcelTable, 0),
		// CharacterSRMap:         make([]*sro.CharacterExcelTable, 0),
		// CharacterSSRMap:        make([]*sro.CharacterExcelTable, 0),
	}

	for _, v := range g.GetExcel().GetCharacterExcelTable() {
		info.CharacterExcelTableMap[v.Id] = v
		if v.IsPlayable && !v.IsNPC && v.IsPlayableCharacter &&
			v.ProductionStep == "Release" && v.Id == v.CharacterPieceItemId {
			info.CharacterReleaseList = append(info.CharacterReleaseList, v)
			info.CharacterMap[v.Id] = v
			// switch v.Rarity {
			// case "R":
			// 	info.CharacterRMap = append(info.CharacterRMap, v)
			// case "SR":
			// 	info.CharacterSRMap = append(info.CharacterSRMap, v)
			// case "SSR":
			// 	info.CharacterSSRMap = append(info.CharacterSSRMap, v)
			// default:
			// 	logger.Debug("未知的角色星级|角色id:%v Rarity:%s", v.Id, v.Rarity)
			// }
		}
	}

	g.GetGPP().CharacterExcel = info
	logger.Info("处理角色完成,角色数量:%v个",
		len(info.CharacterReleaseList))
}

func GetCharacterExcel(characterId int64) *sro.CharacterExcelTable {
	if g := GC.GetGPP(); g == nil {
		return nil
	} else {
		return g.CharacterExcel.CharacterExcelTableMap[characterId]
	}
}

func GetCharacterExcelStruct() *CharacterExcel {
	return GC.GetGPP().CharacterExcel
}

func GetCharacterMap() map[int64]*sro.CharacterExcelTable {
	return GC.GetGPP().CharacterExcel.CharacterMap
}

func RandCharacter() int64 {
	list := GC.GetGPP().CharacterExcel.CharacterReleaseList
	if len(list) == 0 {
		return 0
	}
	return list[rand.Intn(len(list))].Id
}
