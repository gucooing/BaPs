package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
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

func GetCategoryMissionInfo(s *enter.Session) map[string]*sro.CategoryMissionInfo {
	bin := GetMissionBin(s)
	if bin == nil {
		return nil
	}
	if bin.CategoryMissionInfo == nil {
		bin.CategoryMissionInfo = make(map[string]*sro.CategoryMissionInfo)
	}
	return bin.CategoryMissionInfo
}

// 获取每日任务信息
func GetDayMissionInfo(s *enter.Session) *sro.CategoryMissionInfo {
	bin := GetCategoryMissionInfo(s)
	if bin == nil {
		return nil
	}
	if bin["Daily"] == nil {
		bin["Daily"] = &sro.CategoryMissionInfo{}
	}
	info := bin["Daily"]
	// 如果是超了时间就刷新
	if alg.GetDay4().After(time.Unix(info.LastMission, 0)) {
		info.LastMission = time.Now().Unix()
		info.MissionList = make(map[int64]*sro.MissionInfo)
		for _, conf := range gdconf.GetMissionExcelTableCategoryList("Daily") {
			start, _ := time.Parse("2006-01-02 15:04:05", conf.StartDate)
			end, _ := time.Parse("2006-01-02 15:04:05", conf.EndDate)
			if time.Now().After(start) && (time.Now().Before(end) || conf.EndDate == "") {
				info.MissionList[conf.Id] = &sro.MissionInfo{
					MissionId: conf.Id,
					StartTime: time.Now().Unix(),
					Complete:  false,
					Finish:    false,
					ProgressParameters: map[int64]int64{
						proto.GetMissionCompleteConditionType(conf.CompleteConditionType).Value(): 0,
					},
				}
			}
		}
		s.AddMissionByCompleteConditionType(info)
		s.FinishMission(proto.MissionCompleteConditionType_Reset_DailyLogin, 1, nil)
		s.FinishMission(proto.MissionCompleteConditionType_Reset_DailyLoginCount, 1, nil)
	}
	return info
}

// GetMissionInfo 通过任务id拉取已接取的任务
func GetMissionInfo(s *enter.Session, missionId int64) *sro.MissionInfo {
	conf := gdconf.GetMissionExcelTable(missionId)
	if conf == nil {
		return nil
	}
	var bin *sro.CategoryMissionInfo
	switch conf.Category {
	case "Daily":
		bin = GetDayMissionInfo(s)
	default:
		logger.Warn("尚未实现的任务类型:%s", conf.Category)
	}
	if bin == nil {
		logger.Warn("尚未实现的任务类型:%s", conf.Category)
		return nil
	}
	if bin.MissionList == nil {
		bin.MissionList = make(map[int64]*sro.MissionInfo)
	}
	return bin.MissionList[missionId]
}

func GetMissionProgressDB(s *enter.Session, missionId int64) *proto.MissionProgressDB {
	bin := GetMissionInfo(s, missionId)
	if bin == nil {
		return nil
	}
	info := &proto.MissionProgressDB{
		MissionUniqueId:    bin.MissionId,
		Complete:           bin.Complete,
		StartTime:          time.Unix(bin.StartTime, 0),
		ProgressParameters: bin.ProgressParameters,
	}
	return info
}

func GetMissionHistoryDB(s *enter.Session, missionId int64) *proto.MissionHistoryDB {
	bin := GetMissionInfo(s, missionId)
	if bin == nil {
		return nil
	}
	info := &proto.MissionHistoryDB{
		MissionUniqueId: bin.MissionId,
		CompleteTime:    mx.Unix(bin.CompleteTime, 0),
		Expired:         bin.Finish,
	}
	return info
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
		AccountLevel:                int64(GetAccountLevel(s)),
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
