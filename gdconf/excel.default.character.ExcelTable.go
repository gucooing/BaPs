package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadDefaultCharacterExcelTable() {
	g.GetExcel().DefaultCharacterExcelTable = make([]*sro.DefaultCharacterExcelTable, 0)
	name := "DefaultCharacterExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().DefaultCharacterExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetDefaultCharacterExcelTable()))
}

func GetDefaultCharacterExcelTable() []*sro.DefaultCharacterExcelTable {
	if e := GC.GetExcel(); e == nil {
		return nil
	} else {
		return e.GetDefaultCharacterExcelTable()
	}
}
