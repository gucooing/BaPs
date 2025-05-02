package gdconf

import (
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadMultiFloorRaidSeasonManageExcel() {
	g.GetExcel().MultiFloorRaidSeasonManageExcel = make([]*sro.MultiFloorRaidSeasonManageExcel, 0)
	name := "MultiFloorRaidSeasonManageExcel.json"
	mx.LoadExcelJson(g.excelDbPath+name, &g.GetExcel().MultiFloorRaidSeasonManageExcel)
}

type MultiFloorRaidSeasonManage struct {
	Cur                           *sro.MultiFloorRaidSeasonManageExcel
	MultiFloorRaidSeasonManageMap map[int64]*sro.MultiFloorRaidSeasonManageExcel
}

func (g *GameConfig) gppMultiFloorRaidSeasonManageExcel() {
	g.GetGPP().MultiFloorRaidSeasonManage = &MultiFloorRaidSeasonManage{
		MultiFloorRaidSeasonManageMap: make(map[int64]*sro.MultiFloorRaidSeasonManageExcel),
	}
	for _, v := range g.GetExcel().GetMultiFloorRaidSeasonManageExcel() {
		g.GetGPP().MultiFloorRaidSeasonManage.MultiFloorRaidSeasonManageMap[v.SeasonId] = v
	}

	logger.Info("处理制约解除决战赛季配置完成,制约解除决战赛季配置:%v个",
		len(g.GetGPP().MultiFloorRaidSeasonManage.MultiFloorRaidSeasonManageMap))
}

func GetMultiFloorRaidSeasonManageExcel(id int64) *sro.MultiFloorRaidSeasonManageExcel {
	return GC.GetGPP().MultiFloorRaidSeasonManage.MultiFloorRaidSeasonManageMap[id]
}

func GetCurMultiFloorRaidSeasonManageExcel() *sro.MultiFloorRaidSeasonManageExcel {
	conf := GC.GetGPP().MultiFloorRaidSeasonManage

	getCur := func() *sro.MultiFloorRaidSeasonManageExcel {
		for _, v := range conf.MultiFloorRaidSeasonManageMap {
			startTime, err := time.Parse("2006-01-02 15:04:05", v.SeasonStartDate)
			endTime, err := time.Parse("2006-01-02 15:04:05", v.SeasonEndDate)
			if err != nil {
				logger.Error("制约解除决战赛季时间格式错误")
				continue
			}
			if time.Now().After(startTime) && time.Now().Before(endTime) {
				return v
			}
		}
		return nil
	}

	if conf.Cur == nil {
		cur := getCur()
		conf.Cur = cur
	}
	if conf.Cur == nil {
		return nil
	}
	endTime, _ := time.Parse("2006-01-02 15:04:05", conf.Cur.SeasonEndDate)
	if time.Now().After(endTime) {
		conf.Cur = getCur()
	}
	return conf.Cur
}
