package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCafeInfoExcelTable() {
	g.Excel.CafeInfoExcelTableInfo = &sro.CafeInfoExcelTableInfo{
		OrigCafeInfoExcelTable: make([]*sro.CafeInfoExcelTable, 0),
		CafeInfoExcelTables:    make(map[int64]*sro.CafeInfoExcelTable),
	}
	name := "CafeInfoExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.Excel.CafeInfoExcelTableInfo.OrigCafeInfoExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	for _, v := range g.Excel.CafeInfoExcelTableInfo.OrigCafeInfoExcelTable {
		g.Excel.CafeInfoExcelTableInfo.CafeInfoExcelTables[v.CafeId] = v
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.Excel.CafeInfoExcelTableInfo.CafeInfoExcelTables))
}

func GetCafeInfoExcelTables() map[int64]*sro.CafeInfoExcelTable {
	if e := GC.GetExcel(); e == nil {
		return nil
	} else {
		return e.GetCafeInfoExcelTableInfo().GetCafeInfoExcelTables()
	}
}
