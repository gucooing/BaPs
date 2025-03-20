package game

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

var MinElapsedRealtime float64 = 10 // 最低战斗时间

func BattleCheck(s *enter.Session, info *proto.BattleSummary) {
	if s == nil || info == nil {
		return
	}
	level := 0
	// 战斗时间验证
	if info.ElapsedRealtime != 0 &&
		info.ElapsedRealtime < MinElapsedRealtime {
		level++
	}
	// 战斗角色验证

	if level >= 1 {
		logger.Info("[UID:%v]玩家判定为战斗违规", s.AccountServerId)
	}
}

func BattleIsAllAlive(battleSummary *proto.BattleSummary) bool {
	if battleSummary == nil {
		return false
	}
	isSu := true
	for _, heroes := range battleSummary.Group01Summary.Heroes {
		if heroes.HPRateAfter == 0 {
			isSu = false
		}
	}
	return isSu
}

func BattleIsClearTimeInSec(battleSummary *proto.BattleSummary, realtime float64) bool {
	if battleSummary == nil {
		return false
	}
	return battleSummary.ElapsedRealtime < realtime
}

func GetParcelInfo(id, amount int64, v proto.ParcelType) *proto.ParcelInfo {
	return &proto.ParcelInfo{
		Key: &proto.ParcelKeyPair{
			Type: v,
			Id:   id,
		},
		Amount: amount,
		Multiplier: &proto.BasisPoint{
			RawValue: 10000,
		},
		Probability: &proto.BasisPoint{
			RawValue: 10000,
		},
	}
}

// ContentSweepWeekDungeon 悬赏通缉/特别依赖 奖励生成
func ContentSweepWeekDungeon(stageId int64, count int64) ([]*ParcelResult, [][]*proto.ParcelInfo) {
	// 发奖励
	parcelResultList := make([]*ParcelResult, 0)
	clearParcels := make([][]*proto.ParcelInfo, 0)
	for i := int64(0); i < count; i++ {
		clearParcel := make([]*proto.ParcelInfo, 0)
		for _, rewardConf := range gdconf.GetWeekDungeonRewardExcelList(stageId) {
			if !rewardConf.IsDisplayed {
				continue
			}
			parcelType := proto.GetParcelTypeValue(rewardConf.RewardParcelType)
			parcelResultList = append(parcelResultList, &ParcelResult{
				ParcelType: parcelType,
				ParcelId:   rewardConf.RewardParcelId,
				Amount:     rewardConf.RewardParcelAmount,
			})
			clearParcel = append(clearParcel, GetParcelInfo(rewardConf.RewardParcelId,
				rewardConf.RewardParcelAmount, parcelType))
		}
		clearParcels = append(clearParcels, clearParcel)
	}
	return parcelResultList, clearParcels
}

func ContentSweepSchoolDungeon(stageId int64, count int64) ([]*ParcelResult, [][]*proto.ParcelInfo) {
	parcelResultList := make([]*ParcelResult, 0)
	clearParcels := make([][]*proto.ParcelInfo, 0)

	conf := gdconf.GetSchoolDungeonStageExcel(stageId)
	if conf == nil {
		return parcelResultList, clearParcels
	}

	for i := int64(0); i < count; i++ {
		clearParcel := make([]*proto.ParcelInfo, 0)
		for _, rewardConf := range gdconf.GetSchoolDungeonRewardExcelList(conf.StageRewardId) {
			if !rewardConf.IsDisplayed ||
				rewardConf.RewardTag != "Default" {
				continue
			}
			parcelType := proto.GetParcelTypeValue(rewardConf.RewardParcelType)
			parcelResultList = append(parcelResultList, &ParcelResult{
				ParcelType: parcelType,
				ParcelId:   rewardConf.RewardParcelId,
				Amount:     rewardConf.RewardParcelAmount,
			})
			clearParcel = append(clearParcel, GetParcelInfo(rewardConf.RewardParcelId,
				rewardConf.RewardParcelAmount, parcelType))
		}
		clearParcels = append(clearParcels, clearParcel)
	}

	return parcelResultList, clearParcels
}
