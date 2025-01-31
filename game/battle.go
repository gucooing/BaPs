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

	if level >= 1 {
		logger.Info("[UID:%v]玩家判定为战斗违规", s.AccountServerId)
	}
}

// ContentSweepWeekDungeon 副本扫荡
func ContentSweepWeekDungeon(stageId int64, count int64) ([]*ParcelResult, [][]*proto.ParcelInfo, []*proto.ParcelInfo) {
	// 扣钱
	conf := gdconf.GetWeekDungeonExcelTable(stageId)
	if conf == nil {
		return nil, nil, nil
	}

	parcelResultList := make([]*ParcelResult, 0)
	bonusParcels := make([]*proto.ParcelInfo, 0)
	if len(conf.StageEnterCostType) == len(conf.StageEnterCostId) &&
		len(conf.StageEnterCostId) == len(conf.StageEnterCostAmount) {
		for index, rewardType := range conf.StageEnterCostType {
			parcelType := proto.ParcelType(proto.ParcelType_value[rewardType])
			parcelResultList = append(parcelResultList, &ParcelResult{
				ParcelType: parcelType,
				ParcelId:   conf.StageEnterCostId[index],
				Amount:     -conf.StageEnterCostAmount[index] * count,
			})
		}
	}

	// 发奖励
	clearParcels := make([][]*proto.ParcelInfo, 0)
	for i := int64(0); i < count; i++ {
		clearParcel := make([]*proto.ParcelInfo, 0)
		for _, rewardConf := range gdconf.GetWeekDungeonRewardExcelList(stageId) {
			if !rewardConf.IsDisplayed {
				continue
			}
			parcelType := proto.ParcelType(proto.ParcelType_value[rewardConf.RewardParcelType])
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
	return parcelResultList, clearParcels, bonusParcels
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

func GetClanAssistSlotDBs(s *enter.Session) []*proto.ClanAssistSlotDB {
	list := make([]*proto.ClanAssistSlotDB, 0)

	return list
}

func GetAssistCharacterDBs(s *enter.Session) []*proto.AssistCharacterDB {
	list := make([]*proto.AssistCharacterDB, 0)

	return list
}
