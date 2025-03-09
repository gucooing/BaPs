package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadGachaElementExcelTable() {
	g.GetExcel().GachaElementExcelTable = make([]*sro.GachaElementExcelTable, 0)
	name := "GachaElementExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().GachaElementExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GachaElementExcelTable))
}

type GachaElementExcel struct {
	GachaElementExcelMap  map[int64]*sro.GachaElementExcelTable
	GachaElementExcelList map[int64][]*sro.GachaElementExcelTable // GroupId
}

func (g *GameConfig) gppGachaElementExcelTable() {
	g.GetGPP().GachaElementExcel = &GachaElementExcel{
		GachaElementExcelMap:  make(map[int64]*sro.GachaElementExcelTable),
		GachaElementExcelList: make(map[int64][]*sro.GachaElementExcelTable),
	}
	for _, v := range g.GetExcel().GetGachaElementExcelTable() {
		g.GetGPP().GachaElementExcel.GachaElementExcelMap[v.Id] = v
		if g.GetGPP().GachaElementExcel.GachaElementExcelList[v.GachaGroupId] == nil {
			g.GetGPP().GachaElementExcel.GachaElementExcelList[v.GachaGroupId] = make([]*sro.GachaElementExcelTable, 0)
		}
		g.GetGPP().GachaElementExcel.GachaElementExcelList[v.GachaGroupId] = append(
			g.GetGPP().GachaElementExcel.GachaElementExcelList[v.GachaGroupId],
			v,
		)
	}

	logger.Info("处理随机组池配置完成,成就:%v个",
		len(g.GetGPP().GachaElementExcel.GachaElementExcelMap))
}

func GetGachaElementExcelTable(id int64) *sro.GachaElementExcelTable {
	return GC.GetGPP().GachaElementExcel.GachaElementExcelMap[id]
}
