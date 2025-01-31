package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
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

// ContentSweepWeekDungeon 悬赏通缉/特别依赖 奖励生成
func ContentSweepWeekDungeon(stageId int64, count int64) ([]*ParcelResult, [][]*proto.ParcelInfo) {
	// 扣钱
	conf := gdconf.GetWeekDungeonExcelTable(stageId)
	if conf == nil {
		return nil, nil
	}

	parcelResultList := make([]*ParcelResult, 0)
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
	return parcelResultList, clearParcels
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

const (
	AssistTermRewardPeriodFromSec = 20
	AssistRewardLimit             = 1000000000
	AssistRentRewardDailyMaxCount = 20
	AssistRentalFeeAmount         = 40000
)

func GetAssistList(s *enter.Session) map[int32]*sro.AssistList {
	bin := GetEchelonBin(s)
	if bin == nil {
		return nil
	}
	if bin.AssistList == nil {
		bin.AssistList = make(map[int32]*sro.AssistList)
	}
	return bin.AssistList
}

func GetClanAssistSlotDBs(s *enter.Session) []*proto.ClanAssistSlotDB {
	list := make([]*proto.ClanAssistSlotDB, 0)
	for _, assist := range GetAssistList(s) {
		if assist.AssistInfoList == nil {
			assist.AssistInfoList = make(map[int64]*sro.AssistInfo)
		}
		for slot, info := range assist.AssistInfoList {
			clanAssistSlotDB := GetClanAssistSlotDB(s, info)
			if clanAssistSlotDB == nil {
				delete(assist.AssistInfoList, slot)
				continue
			}
			list = append(list, clanAssistSlotDB)
		}
	}

	return list
}

func GetClanAssistSlotDB(s *enter.Session, info *sro.AssistInfo) *proto.ClanAssistSlotDB {
	characterInfo := GetCharacterInfo(s, info.CharacterId)
	if characterInfo == nil {
		return nil
	}
	return &proto.ClanAssistSlotDB{
		EchelonType:      proto.EchelonType(info.EchelonType),
		SlotNumber:       info.SlotNumber,
		CharacterDBId:    characterInfo.ServerId,
		DeployDate:       mx.Unix(info.DeployDate, 0),
		TotalRentCount:   info.TotalRentCount,
		CombatStyleIndex: 0,
	}
}

func GetAssistCharacterDBs(s *enter.Session, assistRelation proto.AssistRelation) []*proto.AssistCharacterDB {
	list := make([]*proto.AssistCharacterDB, 0)
	for _, assist := range GetAssistList(s) {
		for slot, info := range assist.AssistInfoList {
			if assist.AssistInfoList == nil {
				assist.AssistInfoList = make(map[int64]*sro.AssistInfo)
			}
			characterInfo := GetCharacterInfo(s, info.CharacterId)
			if characterInfo == nil {
				delete(assist.AssistInfoList, slot)
				continue
			}
			assistCharacterDB := &proto.AssistCharacterDB{
				EchelonType:             proto.EchelonType(info.EchelonType),
				AccountId:               s.AccountServerId,
				AssistRelation:          assistRelation,
				AssistCharacterServerId: characterInfo.ServerId,
				EquipmentDBs:            make([]*proto.EquipmentDB, 0),
				ExSkillLevel:            characterInfo.ExSkillLevel,
				Exp:                     characterInfo.Exp,
				ExtraPassiveSkillLevel:  characterInfo.ExtraPassiveSkillLevel,
				FavorRank:               characterInfo.FavorRank,
				FavorExp:                characterInfo.FavorExp,
				GearDB:                  GetGearDB(s, characterInfo.GearServerId),
				LeaderSkillLevel:        characterInfo.LeaderSkillLevel,
				Level:                   characterInfo.Level,
				NickName:                GetNickname(s),
				PassiveSkillLevel:       characterInfo.PassiveSkillLevel,
				PotentialStats:          characterInfo.PotentialStats,
				PublicSkillLevel:        characterInfo.CommonSkillLevel,
				SlotNumber:              int32(info.SlotNumber),
				StarGrade:               characterInfo.StarGrade,
				Type:                    proto.ParcelType_Character,
				UniqueId:                characterInfo.CharacterId,
				WeaponDB:                GetWeaponDB(s, characterInfo.CharacterId),

				CostumeId:        0,
				CostumeDB:        nil,
				IsMulligan:       false,
				IsTSAInteraction: false,
				CombatStyleIndex: 0,
			}
			for _, serverId := range characterInfo.EquipmentList {
				assistCharacterDB.EquipmentDBs = append(assistCharacterDB.EquipmentDBs,
					GetEquipmentDB(s, serverId))
			}
			list = append(list, assistCharacterDB)
		}
	}

	return list
}
