package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCafeProductionExcelTable() {
	g.GetExcel().CafeProductionExcelTable = make([]*sro.CafeProductionExcelTable, 0)
	name := "CafeProductionExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CafeProductionExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCafeProductionExcelTable()))
}

type CafeProductionExcel struct {
	CafeProductionExcelMap map[int64]map[int32]map[int64]*sro.CafeProductionExcelTable
}

func (g *GameConfig) gppCafeProductionExcelTable() {
	g.GetGPP().CafeProductionExcel = &CafeProductionExcel{
		CafeProductionExcelMap: make(map[int64]map[int32]map[int64]*sro.CafeProductionExcelTable),
	}

	for _, v := range g.GetExcel().GetCafeProductionExcelTable() {
		if g.GetGPP().CafeProductionExcel.CafeProductionExcelMap[v.CafeId] == nil {
			g.GetGPP().CafeProductionExcel.CafeProductionExcelMap[v.CafeId] = make(map[int32]map[int64]*sro.CafeProductionExcelTable)
		}
		if g.GetGPP().CafeProductionExcel.CafeProductionExcelMap[v.CafeId][v.Rank] == nil {
			g.GetGPP().CafeProductionExcel.CafeProductionExcelMap[v.CafeId][v.Rank] = make(map[int64]*sro.CafeProductionExcelTable)
		}
		g.GetGPP().CafeProductionExcel.CafeProductionExcelMap[v.CafeId][v.Rank][v.CafeProductionParcelId] = v
	}

	logger.Info("处理咖啡厅生产配置表完成,数量:%v个",
		len(g.GetGPP().CafeProductionExcel.CafeProductionExcelMap))
}

func GetCafeProductionExcelTableList(cafeId int64, rank int32) map[int64]*sro.CafeProductionExcelTable {
	if GC.GetGPP().CafeProductionExcel.CafeProductionExcelMap[cafeId] == nil {
		return nil
	}
	return GC.GetGPP().CafeProductionExcel.CafeProductionExcelMap[cafeId][rank]
}
