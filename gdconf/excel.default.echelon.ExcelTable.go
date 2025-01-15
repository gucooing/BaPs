package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadDefaultEchelonExcelTable() {
	g.GetExcel().DefaultEchelonExcelTable = make([]*sro.DefaultEchelonExcelTable, 0)
	name := "DefaultEchelonExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().DefaultEchelonExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().DefaultEchelonExcelTable))
}

func GetDefaultEchelonExcelList() []*sro.DefaultEchelonExcelTable {
	return GC.GetExcel().GetDefaultEchelonExcelTable()
}
