//go:build dev

package gdconf

import (
	"fmt"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
)

func (g *GameConfig) LoadExcel() {
	// 验证文件夹是否存在
	g.excelPath = g.resPath + "/Excel"
	if dirInfo, err := os.Stat(g.excelPath); err != nil || !dirInfo.IsDir() {
		info := fmt.Sprintf("找不到文件夹:%s,err:%s", g.excelPath, err)
		panic(info)
	}
	g.excelPath += "/"

	g.excelDbPath = g.resPath + "/ExcelDB"
	if dirInfo, err := os.Stat(g.excelDbPath); err != nil || !dirInfo.IsDir() {
		info := fmt.Sprintf("找不到文件夹:%s,err:%s", g.excelDbPath, err)
		panic(info)
	}
	g.excelDbPath += "/"

	// 初始化excel
	g.Excel = new(sro.Excel)
	g.loadFunc = []func(){
		g.loadCafeInfoExcelTable,
		g.loadDefaultCharacterExcelTable,
		g.loadCharacterExcelTable,
		g.loadDefaultEchelonExcelTable,
		g.loadDefaultFurnitureExcelTable,
		g.loadShopExcelTable,
		g.loadShopInfoExcelTable,
		g.loadItemExcelTable,
		g.loadEmblemExcel,
		g.loadAcademyFavorScheduleExcelTable,
		g.loadAcademyMessangerExcelTable,
		g.loadGuideMissionExcelTable,
		g.loadScenarioModeExcel,
		g.loadScenarioModeRewardExcel,
		g.loadCharacterWeaponExcelTable,
		g.loadCharacterSkillListExcelTable,
		g.loadSkillExcelTable,
		g.loadRecipeIngredientExcelTable,
		g.loadCampaignStageExcelTable,
		g.loadCampaignUnitExcelTable,
		g.loadWeekDungeonExcelTable,
		g.loadSchoolDungeonStageExcel,
		g.loadAccountLevelExcel,
	}

	for _, fn := range g.loadFunc {
		fn()
	}
}
