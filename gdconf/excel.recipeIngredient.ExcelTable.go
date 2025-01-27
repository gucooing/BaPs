package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadRecipeIngredientExcelTable() {
	g.GetExcel().RecipeIngredientExcelTable = make([]*sro.RecipeIngredientExcelTable, 0)
	name := "RecipeIngredientExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().RecipeIngredientExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetRecipeIngredientExcelTable()))
}

type RecipeIngredientExcel struct {
	RecipeIngredientExcelMap map[int64]*sro.RecipeIngredientExcelTable
}

func (g *GameConfig) gppRecipeIngredientExcelTable() {
	g.GetGPP().RecipeIngredientExcel = &RecipeIngredientExcel{
		RecipeIngredientExcelMap: make(map[int64]*sro.RecipeIngredientExcelTable),
	}

	for _, v := range g.GetExcel().GetRecipeIngredientExcelTable() {
		g.GetGPP().RecipeIngredientExcel.RecipeIngredientExcelMap[v.Id] = v
	}

	logger.Info("处理材料配置表完成,技能配置:%v个",
		len(g.GetGPP().RecipeIngredientExcel.RecipeIngredientExcelMap))
}

func GetRecipeIngredientExcelTable(id int64) *sro.RecipeIngredientExcelTable {
	return GC.GetGPP().RecipeIngredientExcel.RecipeIngredientExcelMap[id]
}
