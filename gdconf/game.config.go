package gdconf

import (
	"runtime"
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

var GC *GameConfig

type GameConfig struct {
	dataPath  string
	resPath   string
	excelPath string
	loadFunc  []func()
	Excel     *sro.Excel
}

func LoadGameConfig(dataPath string, resPath string) *GameConfig {
	gc := new(GameConfig)
	GC = gc
	gc.dataPath = dataPath
	gc.resPath = resPath
	logger.Info("开始读取资源文件")
	startTime := time.Now().Unix()
	gc.LoadExcel()
	endTime := time.Now().Unix()
	runtime.GC()
	logger.Info("读取资源完成,用时:%v秒", endTime-startTime)
	return gc
}

func (g *GameConfig) GetExcel() *sro.Excel {
	if g == nil {
		return nil
	}
	return g.Excel
}
