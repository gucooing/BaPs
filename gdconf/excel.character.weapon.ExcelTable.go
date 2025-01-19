package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCharacterWeaponExcelTable() {
	g.GetExcel().CharacterWeaponExcelTable = make([]*sro.CharacterWeaponExcelTable, 0)
	name := "CharacterWeaponExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CharacterWeaponExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCharacterWeaponExcelTable()))
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
