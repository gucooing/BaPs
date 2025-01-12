package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCharacterExcelTable() {
	g.Excel.CharacterExcelTable = &sro.CharacterExcelTable{
		OrigCharacterExcelTable: make([]*sro.CharacterExcel, 0),
		CharacterExcelMap:       make(map[int64]*sro.CharacterExcel),
	}
	name := "CharacterExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.Excel.CharacterExcelTable.OrigCharacterExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	for _, v := range g.Excel.CharacterExcelTable.OrigCharacterExcelTable {
		g.Excel.CharacterExcelTable.CharacterExcelMap[v.Id] = v
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.Excel.CharacterExcelTable.OrigCharacterExcelTable))
}

func GetCharacterExcel(characterId int64) *sro.CharacterExcel {
	if e := GC.GetExcel(); e == nil {
		return nil
	} else {
		list := e.GetCharacterExcelTable().GetCharacterExcelMap()
		if list == nil {
			return nil
		}
		return list[characterId]
	}
}
