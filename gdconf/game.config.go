package gdconf

import (
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
	CharacterExcel            *CharacterExcel
	CafeInfoExcel             *CafeInfoExcel
	ShopExcel                 *ShopExcel
	ShopInfoExcel             *ShopInfoExcel
	ItemExcel                 *ItemExcel
	Emblem                    *Emblem
	AcademyFavorScheduleExcel *AcademyFavorScheduleExcel
	AcademyMessangerExcel     *AcademyMessangerExcel
	GuideMissionExcel         *GuideMissionExcel
	ScenarioMode              *ScenarioMode
	ScenarioModeReward        *ScenarioModeReward
	CharacterWeaponExcel      *CharacterWeaponExcel
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
