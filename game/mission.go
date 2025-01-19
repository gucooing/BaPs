package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewMission() *sro.MissionBin {
	return &sro.MissionBin{}
}

func GetMissionBin(s *enter.Session) *sro.MissionBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.MissionBin == nil {
		bin.MissionBin = NewMission()
	}
	return bin.MissionBin
}

func GetTutorialList(s *enter.Session) map[int64]bool {
	bin := GetMissionBin(s)
	if bin == nil {
		return nil
	}
	return bin.GetTutorialList()
}

func FinishTutorial(s *enter.Session, tutorialIdS []int64) bool {
	bin := GetMissionBin(s)
	if bin == nil {
		return false
	}
	if bin.TutorialList == nil {
		bin.TutorialList = make(map[int64]bool)
	}
	for _, t := range tutorialIdS {
		bin.TutorialList[t] = true
	}
	return true
}

func GetScenarioGroupHistoryInfoList(s *enter.Session) map[int64]*sro.ScenarioGroupHistoryInfo {
	bin := GetMissionBin(s)
	if bin == nil {
		return nil
	}
	return bin.GetScenarioGroupHistoryInfoList()
}

func GetScenarioGroupHistoryInfo(s *enter.Session, scenarioGroupUniqueId int64) *sro.ScenarioGroupHistoryInfo {
	bin := GetScenarioGroupHistoryInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[scenarioGroupUniqueId]
}

func FinishScenarioGroupHistoryInfo(s *enter.Session, scenarioGroupUniqueId, scenarioType int64) {
	bin := GetMissionBin(s)
	if bin == nil {
		return
	}
	if bin.ScenarioGroupHistoryInfoList == nil {
		bin.ScenarioGroupHistoryInfoList = make(map[int64]*sro.ScenarioGroupHistoryInfo)
	}
	bin.ScenarioGroupHistoryInfoList[scenarioGroupUniqueId] = &sro.ScenarioGroupHistoryInfo{
		ScenarioGroupUqniueId: scenarioGroupUniqueId,
		ClearDateTime:         time.Now().Unix(),
		ScenarioType:          scenarioType,
		// EventContentId:        v.EventContentId,
	}
}

func GetScenarioHistoryInfoList(s *enter.Session) map[int64]*sro.ScenarioHistoryInfo {
	bin := GetMissionBin(s)
	if bin == nil {
		return nil
	}
	return bin.GetScenarioHistoryInfoList()
}

func GetScenarioHistoryInfo(s *enter.Session, scenarioUniqueId int64) *sro.ScenarioHistoryInfo {
	bin := GetScenarioHistoryInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[scenarioUniqueId]
}

func FinishScenarioHistoryInfo(s *enter.Session, scenarioGroupUniqueId int64) []*ParcelResult {
	bin := GetMissionBin(s)
	if bin == nil {
		return nil
	}
	conf := gdconf.GetScenarioModeExcel(scenarioGroupUniqueId)
	if conf == nil {
		return nil
	}
	if bin.ScenarioHistoryInfoList == nil {
		bin.ScenarioHistoryInfoList = make(map[int64]*sro.ScenarioHistoryInfo)
	}
	if _, ok := bin.ScenarioHistoryInfoList[scenarioGroupUniqueId]; ok {
		return nil
	}
	bin.ScenarioHistoryInfoList[scenarioGroupUniqueId] = &sro.ScenarioHistoryInfo{
		ScenarioUniqueId: scenarioGroupUniqueId,
		ClearDateTime:    time.Now().Unix(),
	}
	list := make([]*ParcelResult, 0)
	for _, rewardConf := range gdconf.GetScenarioModeRewardExcel(scenarioGroupUniqueId) {
		list = append(list, &ParcelResult{
			ParcelType: proto.ParcelType(proto.ParcelType_value[rewardConf.RewardParcelType]),
			ParcelId:   rewardConf.RewardParcelId,
			Amount:     rewardConf.RewardParcelAmount,
		})
	}
	return list
}

func GetProgressDBs(s *enter.Session) []*proto.MissionProgressDB {
	list := make([]*proto.MissionProgressDB, 0)

	list = append(list, &proto.MissionProgressDB{
		MissionUniqueId:    1500,
		Complete:           true,
		StartTime:          time.Now(),
		ProgressParameters: make(map[int64]int64),
	})
	list = append(list, &proto.MissionProgressDB{
		MissionUniqueId:    2300,
		Complete:           false,
		StartTime:          time.Now(),
		ProgressParameters: make(map[int64]int64),
	})

	return list
}

func GetDailySuddenMissionInfoDb(s *enter.Session) *proto.MissionInfo {
	info := &proto.MissionInfo{
		Id:                          1606,
		Category:                    proto.MissionCategory_DailySudden,
		ResetType:                   proto.MissionResetType_Daily,
		ToastDisplayType:            0,
		Description:                 3475540373,
		IsVisible:                   true,
		IsLimited:                   false,
		StartDate:                   time.Date(2024, 4, 24, 4, 0, 0, 0, time.UTC),
		StartableEndDate:            time.Date(9999, 12, 31, 23, 59, 59, 9999999, time.UTC),
		EndDate:                     time.Date(9999, 12, 31, 23, 59, 59, 9999999, time.UTC),
		EndDday:                     0,
		AccountState:                s.AccountState,
		AccountLevel:                1,
		PreMissionIds:               make([]int64, 0),
		NextMissionId:               0,
		SuddenMissionContentTypes:   make([]proto.SuddenMissionContentType, 0),
		CompleteConditionType:       proto.MissionCompleteConditionType_Reset_GetItemWithTagCount,
		CompleteConditionCount:      2,
		CompleteConditionParameters: make([]int64, 0),
		RewardIcon:                  "",
		Rewards:                     make([]*proto.ParcelInfo, 0),
		DateAutoRefer:               0,
		ToastImagePath:              "",
		DisplayOrder:                0,
		HasFollowingMission:         false,
		Shortcuts:                   nil,
		ChallengeStageId:            0,
	}

	return info
}
