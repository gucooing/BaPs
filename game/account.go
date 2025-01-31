package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
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

func SetAccountLevel(s *enter.Session, level int32) {
	bin := GetBaseBin(s)
	if bin == nil {
		return
	}
	if level < 0 {
		return
	}
	bin.Level = level
}

func GetNickname(s *enter.Session) string {
	bin := GetBaseBin(s)
	if bin == nil {
		return "hkrpg-go"
	}
	return bin.GetNickname()
}

func GetComment(s *enter.Session) string {
	bin := GetBaseBin(s)
	if bin == nil {
		return "此服务器是免费的"
	}
	return bin.GetComment()
}

func SetComment(s *enter.Session, comment string) {
	bin := GetBaseBin(s)
	if bin == nil {
		return
	}
	bin.Comment = comment
}

func GetEmblemUniqueId(s *enter.Session) int64 {
	bin := GetBaseBin(s)
	if bin == nil {
		return 0
	}
	if bin.EmblemUniqueId == 0 {
		return 1
	}
	return bin.GetEmblemUniqueId()
}

func SetEmblemUniqueId(s *enter.Session, id int64) bool {
	bin := GetBaseBin(s)
	list := GetEmblemInfoList(s)
	if bin == nil || list == nil {
		return false
	}
	if list[id] == nil {
		return false
	}
	bin.EmblemUniqueId = id
	return true
}

func GetLobbyStudent(s *enter.Session) int64 {
	bin := GetBaseBin(s)
	if bin == nil {
		return 0
	}
	if bin.LobbyStudent == 0 {
		return 13010
	}
	return bin.GetLobbyStudent()
}

func GetCardBackgroundId(s *enter.Session) int64 {
	bin := GetBaseBin(s)
	if bin == nil {
		return 0
	}
	if bin.CardBackgroundId == 0 {
		return 1
	}
	return bin.GetCardBackgroundId()
}

func SetCardBackgroundId(s *enter.Session, id int64) bool {
	bin := GetBaseBin(s)
	if bin == nil || GetCardBackgroundIdInfo(s, id) == nil {
		return false
	}
	bin.CardBackgroundId = id
	return true
}

func GetRepresentCharacterUniqueId(s *enter.Session) int64 {
	bin := GetBaseBin(s)
	if bin == nil {
		return 0
	}
	if bin.RepresentCharacterId == 0 {
		return 13010
	}
	return bin.GetRepresentCharacterId()
}

func SetRepresentCharacterUniqueId(s *enter.Session, characterId int64) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	bin.RepresentCharacterId = characterId
	return true
}

func SetLobbyStudent(s *enter.Session, serverId int64) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	bin.LobbyStudent = GetCharacterInfoByServerId(s, serverId).GetCharacterId()
	return true
}

func GetSearchPermission(s *enter.Session) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	return bin.GetSearchPermission()
}

func SetSearchPermission(s *enter.Session, is bool) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	bin.SearchPermission = is
	return true
}

func GetShowAccountLevel(s *enter.Session) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	return bin.GetShowAccountLevel()
}

func SetShowAccountLevel(s *enter.Session, is bool) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	bin.ShowAccountLevel = is
	return true
}

func GetShowFriendCode(s *enter.Session) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	return bin.GetShowFriendCode()
}

func SetShowFriendCode(s *enter.Session, is bool) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	bin.ShowFriendCode = is
	return true
}

func GetShowRaidRanking(s *enter.Session) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	return bin.GetShowRaidRanking()
}

func SetShowRaidRanking(s *enter.Session, is bool) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	bin.ShowRaidRanking = is
	return true
}

func GetShowArenaRanking(s *enter.Session) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	return bin.GetShowArenaRanking()
}

func SetShowArenaRanking(s *enter.Session, is bool) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	bin.ShowArenaRanking = is
	return true
}

func GetShowEliminateRaidRanking(s *enter.Session) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	return bin.GetShowEliminateRaidRanking()
}

func SetShowEliminateRaidRanking(s *enter.Session, is bool) bool {
	bin := GetBaseBin(s)
	if bin == nil {
		return false
	}
	bin.ShowEliminateRaidRanking = is
	return true
}

func GetAccountExp(s *enter.Session) int64 {
	bin := GetBaseBin(s)
	if bin == nil {
		return 0
	}
	return bin.Exp
}

func AddAccountExp(s *enter.Session, num int64) {
	bin := GetBaseBin(s)
	if bin == nil {
		return
	}
	bin.Exp += num
	newLevel, newExp := gdconf.UpAccountLevel(GetAccountLevel(s),
		GetAccountExp(s))
	if bin.Level < newLevel {
		// 升级设置满级体力
		UpCurrency(s, proto.CurrencyTypes_ActionPoint,
			gdconf.GetAPAutoChargeMax(newLevel))
	}
	bin.Exp = newExp
	bin.Level = newLevel
}

func GetLastConnectTime(s *enter.Session) mx.MxTime {
	bin := GetBaseBin(s)
	if bin == nil {
		return mx.MxTime{}
	}
	return mx.Unix(bin.GetLastConnectTime(), 0)
}

func GetAccountDays(s *enter.Session) int32 {
	bin := GetBaseBin(s)
	if bin == nil {
		return 0
	}
	return bin.Days
}

func GetAccountDB(s *enter.Session) *proto.AccountDB {
	baseBin := GetBaseBin(s)
	if s == nil || baseBin == nil {
		logger.Error("账号数据损坏")
		return nil
	}
	info := &proto.AccountDB{
		ServerId:                   baseBin.GetAccountId(),
		Nickname:                   GetNickname(s),
		Level:                      GetAccountLevel(s),
		Exp:                        GetAccountExp(s),
		LastConnectTime:            GetLastConnectTime(s),
		CreateDate:                 time.Unix(baseBin.GetCreateDate(), 0),
		VIPLevel:                   10,
		State:                      s.AccountState,
		Comment:                    GetComment(s),
		RepresentCharacterServerId: GetCharacterServerId(s, GetLobbyStudent(s)),
		PublisherAccountId:         s.YostarUID,
		RetentionDays:              0,
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
