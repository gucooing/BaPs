package gdconf

import (
	"fmt"
	"os"
	"runtime"
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

var GC *GameConfig

type GameConfig struct {
	dataPath    string
	resPath     string
	excelPath   string
	excelDbPath string
	loadFunc    []func()
	gppFunc     []func()
	Excel       *sro.Excel
	GPP         *GPP
}

type GPP struct {
	CharacterExcel                      *CharacterExcel
	CafeInfoExcel                       *CafeInfoExcel
	ShopExcel                           *ShopExcel
	ShopInfoExcel                       *ShopInfoExcel
	ItemExcel                           *ItemExcel
	Emblem                              *Emblem
	AcademyFavorScheduleExcel           *AcademyFavorScheduleExcel
	AcademyMessangerExcel               *AcademyMessangerExcel
	GuideMissionExcel                   *GuideMissionExcel
	ScenarioMode                        *ScenarioMode
	ScenarioModeReward                  *ScenarioModeReward
	CharacterWeaponExcel                *CharacterWeaponExcel
	CharacterSkillListExcel             *CharacterSkillListExcel
	SkillExcel                          *SkillExcel
	RecipeIngredientExcel               *RecipeIngredientExcel
	CampaignStageExcel                  *CampaignStageExcel
	CampaignUnitExcel                   *CampaignUnitExcel
	WeekDungeonExcel                    *WeekDungeonExcel
	SchoolDungeonStage                  *SchoolDungeonStage
	AccountLevel                        *AccountLevel
	EquipmentExcel                      *EquipmentExcel
	EquipmentStatExcel                  *EquipmentStatExcel
	EquipmentLevelExcel                 *EquipmentLevelExcel
	FurnitureExcel                      *FurnitureExcel
	CafeRankExcel                       *CafeRankExcel
	CafeProductionExcel                 *CafeProductionExcel
	IdCardBackground                    *IdCardBackground
	WeekDungeonRewardExcel              *WeekDungeonRewardExcel
	CharacterLevelExcel                 *CharacterLevelExcel
	CharacterGear                       *CharacterGear
	CharacterPotential                  *CharacterPotential
	CharacterPotentialStat              *CharacterPotentialStat
	AcademyLocationExcel                *AcademyLocationExcel
	AcademyTicketExcel                  *AcademyTicketExcel
	AcademyZoneExcel                    *AcademyZoneExcel
	AcademyLocationRankExcel            *AcademyLocationRankExcel
	AcademyRewardExcel                  *AcademyRewardExcel
	SchoolDungeonReward                 *SchoolDungeonReward
	RaidSchedule                        *RaidSchedule
	RaidSeasonManageExcel               *RaidSeasonManageExcel
	RaidRankingRewardExcel              *RaidRankingRewardExcel
	RaidStageExcel                      *RaidStageExcel
	CharacterStatExcel                  *CharacterStatExcel
	RaidStageRewardExcel                *RaidStageRewardExcel
	MissionExcel                        *MissionExcel
	RaidStageSeasonRewardExcel          *RaidStageSeasonRewardExcel
	Attendance                          *Attendance
	MultiFloorRaidReward                *MultiFloorRaidReward
	MultiFloorRaidStage                 *MultiFloorRaidStage
	MultiFloorRaidSeasonManage          *MultiFloorRaidSeasonManage
	TimeAttackDungeonSeasonManageExcel  *TimeAttackDungeonSeasonManageExcel
	RaidEliminateSchedule               *RaidEliminateSchedule
	EliminateRaidSeasonManageExcel      *EliminateRaidSeasonManageExcel
	EliminateRaidRankingRewardExcel     *EliminateRaidRankingRewardExcel
	EliminateRaidStageExcel             *EliminateRaidStageExcel
	EliminateRaidStageRewardExcel       *EliminateRaidStageRewardExcel
	EliminateRaidStageSeasonRewardExcel *EliminateRaidStageSeasonRewardExcel
}

func LoadGameConfig(dataPath string, resPath string) *GameConfig {
	gc := new(GameConfig)
	GC = gc
	gc.dataPath = dataPath
	gc.resPath = resPath
	logger.Info("开始读取资源文件")
	startTime := time.Now().Unix()
	gc.LoadExcel()
	gc.gpp()
	endTime := time.Now().Unix()
	runtime.GC()
	logger.Info("读取资源完成,用时:%v秒", endTime-startTime)
	return gc
}

func (g *GameConfig) gpp() {
	// 验证文件夹是否存在
	if dirInfo, err := os.Stat(g.dataPath); err != nil || !dirInfo.IsDir() {
		info := fmt.Sprintf("找不到文件夹:%s,err:%s", g.dataPath, err)
		panic(info)
	}
	g.dataPath += "/"
	g.GPP = &GPP{}

	g.gppFunc = []func(){
		g.gppCafeInfoExcelTable,
		g.gppCharacterExcelTable,
		g.gppShopExcelTable,
		g.gppShopInfoExcelTable,
		g.gppItemExcelTable,
		g.gppEmblemExcel,
		g.gppAcademyFavorScheduleExcelTable,
		g.gppAcademyMessangerExcelTable,
		g.gppGuideMissionExcelTable,
		g.gppScenarioModeExcel,
		g.gppScenarioModeRewardExcel,
		g.gppCharacterWeaponExcelTable,
		g.gppCharacterSkillListExcelTable,
		g.gppSkillExcelTable,
		g.gppRecipeIngredientExcelTable,
		g.gppCampaignStageExcelTable,
		g.gppCampaignUnitExcelTable,
		g.gppWeekDungeonExcelTable,
		g.gppSchoolDungeonStageExcel,
		g.gppAccountLevelExcel,
		g.gppEquipmentExcelTable,
		g.gppEquipmentStatExcelTable,
		g.gppEquipmentLevelExcelTable,
		g.gppFurnitureExcelTable,
		g.gppCafeRankExcelTable,
		g.gppCafeProductionExcelTable,
		g.gppIdCardBackgroundExcel,
		g.gppWeekDungeonRewardExcelTable,
		g.gppCharacterLevelExcelTable,
		g.gppCharacterGearExcel,
		g.gppCharacterPotentialExcel,
		g.gppCharacterPotentialStatExcel,
		g.gppAcademyLocationExcelTable,
		g.gppAcademyTicketExcelTable,
		g.gppAcademyZoneExcelTable,
		g.gppAcademyLocationRankExcelTable,
		g.gppAcademyRewardExcelTable,
		g.gppSchoolDungeonRewardExcel,
		g.gppRaidSeasonManageExcelTable,
		g.gppRaidRankingRewardExcelTable,
		g.gppRaidStageExcelTable,
		g.gppCharacterStatExcelTable,
		g.gppRaidStageRewardExcelTable,
		g.gppMissionExcelTable,
		g.gppRaidStageSeasonRewardExcelTable,
		g.gppMultiFloorRaidRewardExcel,
		g.gppMultiFloorRaidStageExcel,
		g.gppMultiFloorRaidSeasonManageExcel,
		g.gppTimeAttackDungeonSeasonManageExcelTable,
		g.gppEliminateRaidSeasonManageExcelTable,
		g.gppEliminateRaidRankingRewardExcelTable,
		g.gppEliminateRaidStageExcelTable,
		g.gppEliminateRaidStageRewardExcelTable,
		g.gppEliminateRaidStageSeasonRewardExcelTable,

		// data
		g.loadRaidSchedule,
		g.loadAttendance,
		g.loadRaidEliminateSchedule,
	}

	for _, fn := range g.gppFunc {
		fn()
	}
}

func (g *GameConfig) GetExcel() *sro.Excel {
	if g == nil {
		return nil
	}
	return g.Excel
}

func (g *GameConfig) GetGPP() *GPP {
	if g == nil {
		return nil
	}
	return g.GPP
}
