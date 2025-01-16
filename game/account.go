package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/logger"
)

func GetBaseBin(s *enter.Session) *sro.BasePlayer {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	return bin.GetBaseBin()
}

func GetAccountLevel(s *enter.Session) int32 {
	bin := GetBaseBin(s)
	if bin == nil {
		return 0
	}
	return bin.GetLevel()
}

func GetNickname(s *enter.Session) string {
	bin := GetBaseBin(s)
	if bin == nil {
		return "hkrpg-go"
	}
	return bin.GetNickname()
}

func GetAccountDB(s *enter.Session) *proto.AccountDB {
	baseBin := GetBaseBin(s)
	if s == nil || baseBin == nil {
		logger.Error("账号数据损坏")
		return nil
	}
	info := &proto.AccountDB{
		ServerId:        baseBin.GetAccountId(),
		Nickname:        GetNickname(s),
		Level:           GetAccountLevel(s),
		LastConnectTime: time.Unix(baseBin.GetLastConnectTime(), 0),
		CreateDate:      time.Unix(baseBin.GetCreateDate(), 0),
		VIPLevel:        10,
		State:           s.AccountState,

		RepresentCharacterServerId: 1,
		PublisherAccountId:         s.YostarUID,
		RetentionDays:              1,
	}

	return info
}

func GetAttendanceBookRewards(s *enter.Session) []*proto.AttendanceBookReward {
	return make([]*proto.AttendanceBookReward, 0)
}

func SetAccountNickname(s *enter.Session, nickname string) bool {
	baseBin := s.PlayerBin.GetBaseBin()
	if baseBin == nil {
		return false
	}
	baseBin.Nickname = nickname
	return true
}

func SetLastConnectTime(s *enter.Session) {
	baseBin := s.PlayerBin.GetBaseBin()
	if baseBin == nil {
		return
	}
	baseBin.LastConnectTime = time.Now().Unix()
}

func GetStaticOpenConditions(s *enter.Session) map[string]int32 {
	return map[string]int32{
		"Shop":                              0,
		"Gacha":                             0,
		"LobbyIllust":                       0,
		"Raid":                              2,
		"Cafe":                              2,
		"Unit_Growth_Skill":                 0,
		"Unit_Growth_LevelUp":               0,
		"Unit_Growth_Transcendence":         0,
		"WeekDungeon":                       2,
		"Arena":                             2,
		"Academy":                           2,
		"Equip":                             0,
		"Item":                              0,
		"Mission":                           0,
		"WeekDungeon_Chase":                 2,
		"__Deprecated_WeekDungeon_FindGift": 0,
		"__Deprecated_WeekDungeon_Blood":    0,
		"Story_Sub":                         0,
		"Story_Replay":                      0,
		"None":                              0,
		"Shop_Gem":                          0,
		"Craft":                             2,
		"Student":                           0,
		"GuideMission":                      0,
		"Clan":                              2,
		"Echelon":                           0,
		"Campaign":                          0,
		"EventContent":                      0,
		"EventStage_1":                      2,
		"EventStage_2":                      0,
		"Talk":                              0,
		"Billing":                           0,
		"Schedule":                          0,
		"Story":                             0,
		"Tactic_Speed":                      2,
		"Cafe_Invite":                       2,
		"Cafe_Invite_2":                     18,
		"EventMiniGame_1":                   2,
		"SchoolDungeon":                     2,
		"TimeAttackDungeon":                 2,
		"ShiftingCraft":                     2,
		"Tactic_Skip":                       2,
		"Mulligan":                          2,
		"EventPermanent":                    2,
		"Main_L_1_2":                        32,
		"Main_L_1_3":                        32,
		"Main_L_1_4":                        32,
		"EliminateRaid":                     2,
		"Cafe_2":                            2,
		"MultiFloorRaid":                    2,
		"StrategySkip":                      2,
		"MinigameDreamMaker":                2,
		"MiniGameDefense":                   2,
	}
}
