package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/common/rank"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/mx"
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
		ranks := rank.NewArenaRank(conf.GetUniqueId(), s.AccountServerId)
		old := bin.ArenaBin
		info := &sro.ArenaBin{
			CurSeasonId:   conf.GetUniqueId(),
			PlayerGroupId: 1, // 默认全放到一起去
			SeasonRecord:  ranks,
			AllTimeRecord: alg.MaxInt64(old.GetAllTimeRecord(), ranks),
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

		CumulativeTimeReward:     0,                              // 积累的时间奖励
		TimeRewardLastUpdateTime: time.Now(),                     // 奖励最后更新时间
		DailyRewardActiveTime:    time.Now().Add(24 * time.Hour), // 下一个每日排名奖励可领取时间
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

	ranks := make(map[int64]bool, 0)
	r := rank.GetArenaRank(bin.GetCurSeasonId(), s.AccountServerId)
	if r <= 3 {
		for i := int64(1); i < 5; i++ {
			if i == r {
				continue
			}
			ranks[i] = true
		}
	} else {
		for i := 0; i < 3; i++ {
			uid := rand.Int63n(r-1) + 1
			if ranks[uid] {
				i--
				continue
			}
			ranks[uid] = true
		}
	}

	list := make([]*proto.ArenaUserDB, 0)
	for aernaRank := range ranks {
		var info *proto.ArenaUserDB
		uid, _ := rank.GetArenaUidByRank(bin.GetCurSeasonId(), aernaRank)
		if ps := enter.GetSessionByUid(uid); ps != nil {
			// 补上真人
		} else {
			info = GetNPCArenaUserDB(aernaRank)
		}

		list = append(list, info)
	}

	return list
}

func GetNPCArenaUserDB(aernaRank int64) *proto.ArenaUserDB {
	conf := gdconf.GetArenaNPCInfo()
	if conf == nil {
		conf = gdconf.DefaultArenaNPCInfo
	}
	characterId := gdconf.RandCharacter()
	info := &proto.ArenaUserDB{
		RepresentCharacterUniqueId: characterId,
		NickName:                   fmt.Sprintf("Character_%v", characterId),
		Rank:                       aernaRank,
		Level:                      conf.NpcaccountLevel,
		TeamSettingDB: &proto.ArenaTeamSettingDB{
			EchelonType:               proto.EchelonType_ArenaDefence,
			LeaderCharacterId:         characterId,
			TSSInteractionCharacterId: 0,
			MainCharacters:            make([]*proto.ArenaCharacterDB, 0),
			SupportCharacters:         make([]*proto.ArenaCharacterDB, 0),
			MapId:                     1006,
		},
	}
	for _, id := range conf.ExceptionMainCharacterIds {
		info.TeamSettingDB.MainCharacters = append(
			info.TeamSettingDB.MainCharacters, GetNPCArenaCharacterDB(id),
		)
	}
	for _, id := range conf.ExceptionSupportCharacterIds {
		info.TeamSettingDB.SupportCharacters = append(
			info.TeamSettingDB.SupportCharacters, GetNPCArenaCharacterDB(id),
		)
	}

	return info
}

func GetNPCArenaCharacterDB(id int64) *proto.ArenaCharacterDB {
	info := &proto.ArenaCharacterDB{
		UniqueId:               id,
		StarGrade:              3,
		Level:                  22,
		PublicSkillLevel:       1,
		ExSkillLevel:           1,
		PassiveSkillLevel:      1,
		ExtraPassiveSkillLevel: 1,
		LeaderSkillLevel:       1,
	}
	return info
}
