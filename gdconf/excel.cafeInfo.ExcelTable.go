package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCafeInfoExcelTable() {
	g.GetExcel().CafeInfoExcelTableInfo = make([]*sro.CafeInfoExcelTableInfo, 0)
	name := "CafeInfoExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CafeInfoExcelTableInfo); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCafeInfoExcelTableInfo()))
}

type CafeInfoExcel struct {
	CafeInfoExcelTableMap map[int64]*sro.CafeInfoExcelTableInfo
}

func (g *GameConfig) gppCafeInfoExcelTable() {
	g.GetGPP().CafeInfoExcel = &CafeInfoExcel{
		CafeInfoExcelTableMap: make(map[int64]*sro.CafeInfoExcelTableInfo, 0),
	}

	for _, v := range g.GetExcel().GetCafeInfoExcelTableInfo() {
		g.GetGPP().CafeInfoExcel.CafeInfoExcelTableMap[v.CafeId] = v
	}

	logger.Info("处理咖啡厅数量完成:%v个", len(g.GetGPP().CafeInfoExcel.CafeInfoExcelTableMap))
}

func GetCafeInfoExcelTables() map[int64]*sro.CafeInfoExcelTableInfo {
	if g := GC.GetGPP(); g == nil {
		return nil
	} else {
		return g.CafeInfoExcel.CafeInfoExcelTableMap
	}
}
