package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCharacterExcelTable() {
	g.GetExcel().CharacterExcelTable = make([]*sro.CharacterExcelTable, 0)
	name := "CharacterExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CharacterExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCharacterExcelTable()))
}

type CharacterExcel struct {
	CharacterExcelTableMap map[int64]*sro.CharacterExcelTable
	CharacterMap           map[int64]*sro.CharacterExcelTable // 全部角色索引
	CharacterRMap          []*sro.CharacterExcelTable
	CharacterSRMap         []*sro.CharacterExcelTable
	CharacterSSRMap        []*sro.CharacterExcelTable
}

/*
R 1
SR 2
SSR 3
*/

func (g *GameConfig) gppCharacterExcelTable() {
	info := &CharacterExcel{
		CharacterExcelTableMap: make(map[int64]*sro.CharacterExcelTable),
		CharacterMap:           make(map[int64]*sro.CharacterExcelTable),
		CharacterRMap:          make([]*sro.CharacterExcelTable, 0),
		CharacterSRMap:         make([]*sro.CharacterExcelTable, 0),
		CharacterSSRMap:        make([]*sro.CharacterExcelTable, 0),
	}

	for _, v := range g.GetExcel().GetCharacterExcelTable() {
		info.CharacterExcelTableMap[v.Id] = v
		if v.IsPlayable && !v.IsNPC && v.IsPlayableCharacter &&
			v.ProductionStep_ == "Release" {
			info.CharacterMap[v.Id] = v
			switch v.Rarity_ {
			case "R":
				info.CharacterRMap = append(info.CharacterRMap, v)
			case "SR":
				info.CharacterSRMap = append(info.CharacterSRMap, v)
			case "SSR":
				info.CharacterSSRMap = append(info.CharacterSSRMap, v)
			default:
				logger.Debug("未知的角色星级|角色id:%v Rarity:%s", v.Id, v.Rarity_)
			}
		}
	}

	g.GetGPP().CharacterExcel = info
	logger.Info("处理角色完成,三星角色:%v个,二星角色:%v个,一星角色:%v个",
		len(info.CharacterSSRMap), len(info.CharacterSRMap), len(info.CharacterRMap))
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
