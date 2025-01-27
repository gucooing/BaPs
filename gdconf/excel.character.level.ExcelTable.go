package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCharacterLevelExcelTable() {
	g.GetExcel().CharacterLevelExcelTable = make([]*sro.CharacterLevelExcelTable, 0)
	name := "CharacterLevelExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CharacterLevelExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCharacterLevelExcelTable()))
}

type CharacterLevelExcel struct {
	CharacterLevelExcelTableMap map[int32]*sro.CharacterLevelExcelTable
}

func (g *GameConfig) gppCharacterLevelExcelTable() {
	g.GetGPP().CharacterLevelExcel = &CharacterLevelExcel{
		CharacterLevelExcelTableMap: make(map[int32]*sro.CharacterLevelExcelTable, 0),
	}
	for _, v := range g.GetExcel().GetCharacterLevelExcelTable() {
		g.GetGPP().CharacterLevelExcel.CharacterLevelExcelTableMap[v.Level] = v
	}
	logger.Info("处理角色等级配置表完成数量:%v个", len(g.GetGPP().CharacterLevelExcel.CharacterLevelExcelTableMap))
}

func GetCharacterLevelExcelTable(level int32) *sro.CharacterLevelExcelTable {
	return GC.GetGPP().CharacterLevelExcel.CharacterLevelExcelTableMap[level]
}

func UpCharacterLevel(level int32, exp int64) (int32, int64) {
	for {
		conf := GetCharacterLevelExcelTable(level)
		if conf == nil {
			return level - 1, exp
		}
		if exp < conf.Exp {
			return level, exp
		}
		exp -= conf.Exp
		level++
	}
}
