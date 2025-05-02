package game

import (
	"fmt"
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/common/rank"
	ranar "github.com/gucooing/BaPs/common/rank_arena"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
)

const (
	ArenaBattleTime = 60
)

func GetArenaBin(s *enter.Session) *sro.ArenaBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	conf := gdconf.GetCurArenaSeason()
	if conf == nil {
		return nil
	}
	newArena := func() {
		old := bin.ArenaBin
		info := &sro.ArenaBin{
			CurSeasonId:   conf.GetUniqueId(),
			PlayerGroupId: 1, // 默认全放到一起去
			SeasonRecord:  ranar.DefaultArenaRank,
			AllTimeRecord: alg.MaxInt64(old.GetAllTimeRecord(), ranar.DefaultArenaRank),
		}
		bin.ArenaBin = info
	}
	if bin.ArenaBin == nil {
		newArena()
	}
	if conf.UniqueId != bin.ArenaBin.CurSeasonId {
		newArena()
	}
	return bin.ArenaBin
}

func GetArenBattleEnterActiveTime(s *enter.Session) time.Time {
	bin := GetArenaBin(s)
	return time.Unix(bin.GetBattleEnterTime(), 0).Add(ArenaBattleTime * time.Second)
}

func GetArenAutoRefreshTime(s *enter.Session) time.Time {
	bin := GetArenaBin(s)
	return time.Unix(bin.GetAutoRefreshTime(), 0)
}

func GetArenaPlayerInfoDB(s *enter.Session) *proto.ArenaPlayerInfoDB {
	bin := GetArenaBin(s)
	info := &proto.ArenaPlayerInfoDB{
		CurrentSeasonId:       bin.GetCurSeasonId(),                                       // 当前赛季
		PlayerGroupId:         bin.GetPlayerGroupId(),                                     // 玩家组
		CurrentRank:           rank.GetArenaRank(bin.GetCurSeasonId(), s.AccountServerId), // 当前排名
		SeasonRecord:          bin.GetSeasonRecord(),                                      // 本赛季最高记录
		AllTimeRecord:         bin.GetAllTimeRecord(),                                     // 历史记录
		BattleEnterActiveTime: mx.MxTime(GetArenBattleEnterActiveTime(s)),                 // 战斗冷却结束时间

		CumulativeTimeReward:     0,                            // 积累的时间奖励
		TimeRewardLastUpdateTime: mx.Now(),                     // 奖励最后更新时间
		DailyRewardActiveTime:    mx.Now().Add(24 * time.Hour), // 下一个每日排名奖励可领取时间
	}
	return info
}

func GetOpponentUserDBs(s *enter.Session) []*proto.ArenaUserDB {
	bin := GetArenaBin(s)
	// 刷新自动刷新时间
	if bin == nil {
		return nil
	}
	bin.AutoRefreshTime = time.Now().Unix()

	s.GenArenaUserList(bin.GetCurSeasonId())
	list := make([]*proto.ArenaUserDB, 0)
	for _, v := range s.GetArenaUserList() {
		if ps := enter.GetSessionByUid(v.Uid); ps != nil && !v.IsNpc {
			list = append(list, GetPlayerArenaUserDB(ps, proto.EchelonType_ArenaDefence))
		} else {
			list = append(list, GetNPCArenaUserDB(v))
		}
	}

	return list
}

func GetNPCArenaUserDB(au *enter.ArenaUser) *proto.ArenaUserDB {
	conf := gdconf.GetArenaNPCByIndex(au.Index)
	info := &proto.ArenaUserDB{
		RepresentCharacterUniqueId: au.CharacterId,
		NickName:                   fmt.Sprintf("Character_%v", au.CharacterId),
		Rank:                       au.Rank,
		Level:                      conf.NpcaccountLevel,
		TeamSettingDB: &proto.ArenaTeamSettingDB{
			EchelonType:               proto.EchelonType_ArenaDefence,
			LeaderCharacterId:         au.CharacterId,
			TSSInteractionCharacterId: 0,
			MainCharacters:            make([]*proto.ArenaCharacterDB, 0),
			SupportCharacters:         make([]*proto.ArenaCharacterDB, 0),
			MapId:                     1006,
		},
	}
	for _, id := range conf.ExceptionMainCharacterIds {
		info.TeamSettingDB.MainCharacters = append(
			info.TeamSettingDB.MainCharacters, conf.GetArenaCharacterDB(id),
		)
	}
	for _, id := range conf.ExceptionSupportCharacterIds {
		info.TeamSettingDB.SupportCharacters = append(
			info.TeamSettingDB.SupportCharacters, conf.GetArenaCharacterDB(id),
		)
	}

	return info
}

func GetPlayerArenaUserDB(s *enter.Session, echelonType proto.EchelonType) *proto.ArenaUserDB {
	echelonInfo := GetEchelonInfo(s, int32(echelonType), 1)
	info := &proto.ArenaUserDB{
		RepresentCharacterUniqueId: GetRepresentCharacterUniqueId(s),
		NickName:                   GetNickname(s),
		Level:                      int64(GetAccountLevel(s)),
		TeamSettingDB: &proto.ArenaTeamSettingDB{
			EchelonType:               echelonType,
			LeaderCharacterId:         echelonInfo.GetLeaderCharacter(),
			TSSInteractionCharacterId: echelonInfo.GetTssId(),
			MainCharacters:            make([]*proto.ArenaCharacterDB, 0),
			SupportCharacters:         make([]*proto.ArenaCharacterDB, 0),
			MapId:                     1006,
		},
	}
	for _, characterId := range echelonInfo.GetMainCharacterList() {
		ac := GetArenaCharacterDB(s, characterId)
		if ac == nil {
			continue
		}
		info.TeamSettingDB.MainCharacters = append(
			info.TeamSettingDB.MainCharacters, ac,
		)
	}
	for _, characterId := range echelonInfo.GetSupportCharacterList() {
		ac := GetArenaCharacterDB(s, characterId)
		if ac == nil {
			continue
		}
		info.TeamSettingDB.SupportCharacters = append(
			info.TeamSettingDB.SupportCharacters, ac,
		)
	}

	return info
}

func GetArenaCharacterDB(s *enter.Session, characterId int64) *proto.ArenaCharacterDB {
	db := GetCharacterInfo(s, characterId)
	if db == nil {
		return nil
	}
	info := &proto.ArenaCharacterDB{
		ServerId:               db.ServerId,
		UniqueId:               db.CharacterId,
		StarGrade:              db.StarGrade,
		Level:                  db.Level,
		PublicSkillLevel:       db.CommonSkillLevel,
		ExSkillLevel:           db.ExSkillLevel,
		PassiveSkillLevel:      db.PassiveSkillLevel,
		ExtraPassiveSkillLevel: db.ExtraPassiveSkillLevel,
		LeaderSkillLevel:       db.LeaderSkillLevel,
		EquipmentDBs:           make([]*proto.EquipmentDB, 0),
		FavorRankInfo:          make(map[int64]int64),
		PotentialStats:         db.PotentialStats,
		WeaponDB:               GetWeaponDB(s, db.CharacterId),
		GearDB:                 GetGearDB(s, db.GearServerId),

		CostumeDB:        nil,
		CombatStyleIndex: 0,
	}

	return info
}
