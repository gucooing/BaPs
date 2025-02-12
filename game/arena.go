package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetArenaPlayerInfoDB(s *enter.Session) *proto.ArenaPlayerInfoDB {
	info := &proto.ArenaPlayerInfoDB{
		CurrentSeasonId:          1,                              // 当前赛季
		PlayerGroupId:            1,                              // 玩家组
		CurrentRank:              1,                              // 当前排名
		SeasonRecord:             1,                              // 本赛季最高记录
		AllTimeRecord:            0,                              // 历史记录
		CumulativeTimeReward:     0,                              // 积累的时间奖励
		TimeRewardLastUpdateTime: time.Now(),                     // 奖励最后更新时间
		BattleEnterActiveTime:    time.Now(),                     // 战斗冷却结束时间
		DailyRewardActiveTime:    time.Now().Add(24 * time.Hour), // 下一个每日排名奖励可领取时间
	}
	return info
}

func GetOpponentUserDBs(s *enter.Session) []*proto.ArenaUserDB {
	list := make([]*proto.ArenaUserDB, 0)
	for i := int64(0); i < 3; i++ {
		info := &proto.ArenaUserDB{
			AccountServerId:             0,
			RepresentCharacterUniqueId:  10082,
			RepresentCharacterCostumeId: 0,
			NickName:                    "Character_10008",
			Rank:                        i + 2,
			Level:                       25,
			Exp:                         0,
			TeamSettingDB: &proto.ArenaTeamSettingDB{
				EchelonType:               proto.EchelonType_ArenaDefence,
				LeaderCharacterId:         10082,
				TSSInteractionCharacterId: 0,
				MainCharacters: []*proto.ArenaCharacterDB{
					GetArenaCharacterDB(10082),
					GetArenaCharacterDB(10090),
					GetArenaCharacterDB(13003),
					GetArenaCharacterDB(10083),
				},
				SupportCharacters: []*proto.ArenaCharacterDB{
					GetArenaCharacterDB(23003),
					GetArenaCharacterDB(26005),
				},
				TSSCharacterDB: nil,
				MapId:          1004,
			},
			AccountAttachmentDB: nil,
			UserName:            "",
		}
		list = append(list, info)
	}

	return list
}

func GetArenaCharacterDB(id int64) *proto.ArenaCharacterDB {
	info := &proto.ArenaCharacterDB{
		ServerId:               0,
		UniqueId:               id,
		StarGrade:              3,
		Level:                  22,
		PublicSkillLevel:       1,
		ExSkillLevel:           1,
		PassiveSkillLevel:      1,
		ExtraPassiveSkillLevel: 1,
		LeaderSkillLevel:       1,
		EquipmentDBs:           make([]*proto.EquipmentDB, 0),

		FavorRankInfo:    nil,
		PotentialStats:   nil,
		CombatStyleIndex: 0,
		WeaponDB:         nil,
		GearDB:           nil,
		CostumeDB:        nil,
	}
	return info
}
