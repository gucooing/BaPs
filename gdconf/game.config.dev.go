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
		g.loadEquipmentExcelTable,
		g.loadEquipmentStatExcelTable,
		g.loadEquipmentLevelExcelTable,
		g.loadFurnitureExcelTable,
		g.loadCafeRankExcelTable,
		g.loadCafeProductionExcelTable,
		g.loadIdCardBackgroundExcel,
		g.loadWeekDungeonRewardExcelTable,
		g.loadCharacterLevelExcelTable,
		g.loadCharacterGearExcel,
		g.loadCharacterPotentialExcel,
		g.loadCharacterPotentialStatExcel,
		g.loadAcademyLocationExcelTable,
		g.loadAcademyTicketExcelTable,
		g.loadAcademyZoneExcelTable,
		g.loadAcademyLocationRankExcelTable,
		g.loadAcademyRewardExcelTable,
		g.loadSchoolDungeonRewardExcel,
		g.loadRaidSeasonManageExcelTable,
		g.loadRaidRankingRewardExcelTable,
		g.loadRaidStageExcelTable,
		g.loadCharacterStatExcelTable,
		g.loadRaidStageRewardExcelTable,
		g.loadMissionExcelTable,
		g.loadRaidStageSeasonRewardExcelTable,
		g.loadMultiFloorRaidRewardExcel,
		g.loadMultiFloorRaidStageExcel,
		g.loadMultiFloorRaidSeasonManageExcel,
		g.loadTimeAttackDungeonSeasonManageExcelTable,
		g.loadEliminateRaidSeasonManageExcelTable,
		g.loadEliminateRaidRankingRewardExcelTable,
		g.loadEliminateRaidStageExcelTable,
		g.loadEliminateRaidStageRewardExcelTable,
		g.loadEliminateRaidStageSeasonRewardExcelTable,
		g.loadGachaElementRecursiveExcelTable,
		g.loadGachaElementExcelTable,
		g.loadGoodsExcelTable,
		g.loadTimeAttackDungeonGeasExcelTable,
		g.loadTimeAttackDungeonRewardExcelTable,
		g.loadShopRefreshExcelTable,
		g.loadArenaSeasonExcelTable,
	}

	for _, fn := range g.loadFunc {
		fn()
	}
}
