package enter

import (
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

type Mission struct {
	MissionByCompleteConditionType map[string]map[int64]*sro.MissionInfo
	SyncMission                    map[int64]*sro.MissionInfo // 待同步的任务列表
}

func (x *Session) GetMission() *Mission {
	if x == nil {
		return nil
	}
	if x.Mission == nil {
		x.Mission = &Mission{}
	}
	return x.Mission
}

func (x *Session) AddMissionByCompleteConditionType(info *sro.CategoryMissionInfo) {
	bin := x.GetMission()
	if bin == nil {
		return
	}
	if bin.MissionByCompleteConditionType == nil {
		bin.MissionByCompleteConditionType = make(map[string]map[int64]*sro.MissionInfo)
	}
	for id, mission := range info.GetMissionList() {
		conf := gdconf.GetMissionExcelTable(id)
		if conf == nil {
			continue
		}
		if bin.MissionByCompleteConditionType[conf.CompleteConditionType] == nil {
			bin.MissionByCompleteConditionType[conf.CompleteConditionType] = make(map[int64]*sro.MissionInfo)
		}
		bin.MissionByCompleteConditionType[conf.CompleteConditionType][id] = mission
	}
}

func (x *Session) GetMissionByCompleteConditionType() map[string]map[int64]*sro.MissionInfo {
	bin := x.GetMission()
	if bin == nil {
		return nil
	}
	if bin.MissionByCompleteConditionType == nil {
		bin.MissionByCompleteConditionType = make(map[string]map[int64]*sro.MissionInfo)
	}
	return bin.MissionByCompleteConditionType
}

func (x *Session) GetMissionSync() map[int64]*sro.MissionInfo {
	bin := x.GetMission()
	if bin == nil {
		return nil
	}
	if bin.SyncMission == nil {
		bin.SyncMission = make(map[int64]*sro.MissionInfo)
	}
	return bin.SyncMission
}

func (x *Mission) AddMissionSync(info *sro.MissionInfo) {
	if x == nil || info == nil {
		return
	}
	if x.SyncMission == nil {
		x.SyncMission = make(map[int64]*sro.MissionInfo)
	}
	x.SyncMission[info.MissionId] = info
}

func (x *Session) NewMissionSync() {
	bin := x.GetMission()
	if bin == nil {
		return
	}
	bin.SyncMission = make(map[int64]*sro.MissionInfo)
}

func (x *Session) MissionReward(info *sro.MissionInfo) {
	if x == nil || info == nil {
		return
	}
	info.CompleteTime = time.Now().Unix()
	info.Finish = true
	x.GetMission().AddMissionSync(info)
}

func (x *Session) FinishMission(t proto.MissionCompleteConditionType, num int64, parameter []int64) {
	binList := x.GetMissionByCompleteConditionType()[t.String()]
	for _, bin := range binList {
		if bin.Complete {
			continue
		}
		switch t {
		// 计数类型的完成条件
		case proto.MissionCompleteConditionType_Reset_DailyLogin,
			proto.MissionCompleteConditionType_Reset_DailyLoginCount,
			proto.MissionCompleteConditionType_Reset_LoginAtSpecificTime:
			x.GetMission().MissionFinishNum(t, bin, num)
		default:
			logger.Warn("未知的任务完成类型:%s", t.String())
		}
	}
}

// 计数类型的完成条件
func (x *Mission) MissionFinishNum(t proto.MissionCompleteConditionType, bin *sro.MissionInfo, num int64) {
	if bin == nil {
		return
	}
	conf := gdconf.GetMissionExcelTable(bin.MissionId)
	if conf == nil || bin.Complete { // 没有配置和已完成的排除掉
		return
	}
	if bin.ProgressParameters == nil {
		bin.ProgressParameters = map[int64]int64{
			int64(t): 0,
		}
	}
	old := bin.ProgressParameters[int64(t)]
	bin.ProgressParameters[int64(t)] = alg.MaxInt64(bin.ProgressParameters[int64(t)]+num,
		conf.CompleteConditionCount)
	if bin.ProgressParameters[int64(t)] == conf.CompleteConditionCount {
		bin.Complete = true
	}
	if bin.ProgressParameters[int64(t)] != old {
		x.AddMissionSync(bin)
	}
}
