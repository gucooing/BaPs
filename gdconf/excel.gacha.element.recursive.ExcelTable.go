package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadGachaElementRecursiveExcelTable() {
	g.GetExcel().GachaElementRecursiveExcelTable = make([]*sro.GachaElementRecursiveExcelTable, 0)
	name := "GachaElementRecursiveExcelTable.json"
	loadExcelJson(g.excelPath+name, &g.GetExcel().GachaElementRecursiveExcelTable)
}

type GachaElementRecursiveExcel struct {
	GachaElementRecursiveExcelMap  map[int64]*sro.GachaElementRecursiveExcelTable
	GachaElementRecursiveExcelList map[int64][]*sro.GachaElementRecursiveExcelTable // GroupId
}

func (g *GameConfig) gppGachaElementRecursiveExcelTable() {
	g.GetGPP().GachaElementRecursiveExcel = &GachaElementRecursiveExcel{
		GachaElementRecursiveExcelMap:  make(map[int64]*sro.GachaElementRecursiveExcelTable),
		GachaElementRecursiveExcelList: make(map[int64][]*sro.GachaElementRecursiveExcelTable),
	}
	for _, v := range g.GetExcel().GetGachaElementRecursiveExcelTable() {
		g.GetGPP().GachaElementRecursiveExcel.GachaElementRecursiveExcelMap[v.Id] = v
		if g.GetGPP().GachaElementRecursiveExcel.GachaElementRecursiveExcelList[v.GachaGroupId] == nil {
			g.GetGPP().GachaElementRecursiveExcel.GachaElementRecursiveExcelList[v.GachaGroupId] = make([]*sro.GachaElementRecursiveExcelTable, 0)
		}
		g.GetGPP().GachaElementRecursiveExcel.GachaElementRecursiveExcelList[v.GachaGroupId] = append(
			g.GetGPP().GachaElementRecursiveExcel.GachaElementRecursiveExcelList[v.GachaGroupId],
			v,
		)
	}

	logger.Info("处理随机组配置完成,成就:%v个",
		len(g.GetGPP().GachaElementRecursiveExcel.GachaElementRecursiveExcelMap))
}

func GetGachaElementRecursiveExcelTable(id int64) *sro.GachaElementRecursiveExcelTable {
	return GC.GetGPP().GachaElementRecursiveExcel.GachaElementRecursiveExcelMap[id]
}

func GetGachaElementExcelTableByGachaGroupId(gachaGroupId int64) []*GachaElementGroupId {
	list := make([]*GachaElementGroupId, 0)
	recursiveList := GC.GetGPP().GachaElementRecursiveExcel.GachaElementRecursiveExcelList[gachaGroupId]
	if len(recursiveList) == 0 {
		list = append(list, GetGachaElementGroupIdByGachaGroupId(gachaGroupId))
	} else {
		for _, info := range recursiveList {
			list = append(list, GetGachaElementGroupIdByGachaGroupId(info.ParcelId))
			// conf := GetGachaElementGroupIdByGachaGroupId(info.ParcelId)
		}
	}
	return list
}
