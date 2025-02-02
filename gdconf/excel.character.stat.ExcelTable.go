package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCharacterStatExcelTable() {
	g.GetExcel().CharacterStatExcelTable = make([]*sro.CharacterStatExcelTable, 0)
	name := "CharacterStatExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CharacterStatExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCharacterStatExcelTable()))
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
