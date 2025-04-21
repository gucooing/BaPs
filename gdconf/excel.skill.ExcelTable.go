package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadSkillExcelTable() {
	g.GetExcel().SkillExcelTable = make([]*sro.SkillExcelTable, 0)
	name := "SkillExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().SkillExcelTable)
}

type SkillExcel struct {
	SkillExcelMap map[string]map[int32]*sro.SkillExcelTable
}

func (g *GameConfig) gppSkillExcelTable() {
	g.GetGPP().SkillExcel = &SkillExcel{
		SkillExcelMap: make(map[string]map[int32]*sro.SkillExcelTable),
	}

	for _, v := range g.GetExcel().GetSkillExcelTable() {
		if g.GetGPP().SkillExcel.SkillExcelMap[v.GroupId] == nil {
			g.GetGPP().SkillExcel.SkillExcelMap[v.GroupId] = make(map[int32]*sro.SkillExcelTable)
		}
		g.GetGPP().SkillExcel.SkillExcelMap[v.GroupId][v.Level] = v
	}

	logger.Info("处理技能配置完成,技能配置:%v个",
		len(g.GetGPP().SkillExcel.SkillExcelMap))
}

func GetSkillExcelTable(groupId string, level int32) *sro.SkillExcelTable {
	list, ok := GC.GetGPP().SkillExcel.SkillExcelMap[groupId]
	if !ok {
		return nil
	}
	return list[level]
}
