package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCafeRankExcelTable() {
	g.GetExcel().CafeRankExcelTable = make([]*sro.CafeRankExcelTable, 0)
	name := "CafeRankExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CafeRankExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCafeRankExcelTable()))
}

type CafeRankExcel struct {
	CafeRankExcelMap map[int64]map[int32]*sro.CafeRankExcelTable
}

func (g *GameConfig) gppCafeRankExcelTable() {
	g.GetGPP().CafeRankExcel = &CafeRankExcel{
		CafeRankExcelMap: make(map[int64]map[int32]*sro.CafeRankExcelTable),
	}

	for _, v := range g.GetExcel().GetCafeRankExcelTable() {
		if g.GetGPP().CafeRankExcel.CafeRankExcelMap[v.CafeId] == nil {
			g.GetGPP().CafeRankExcel.CafeRankExcelMap[v.CafeId] = make(map[int32]*sro.CafeRankExcelTable)
		}
		g.GetGPP().CafeRankExcel.CafeRankExcelMap[v.CafeId][v.Rank] = v
	}

	logger.Info("处理咖啡厅等级配置表完成,数量:%v个",
		len(g.GetGPP().CafeRankExcel.CafeRankExcelMap))
}

func GetCafeRankExcelTable(cafeId int64, rank int32) *sro.CafeRankExcelTable {
	if GC.GetGPP().CafeRankExcel.CafeRankExcelMap[cafeId] == nil {
		return nil
	}
	return GC.GetGPP().CafeRankExcel.CafeRankExcelMap[cafeId][rank]
}
