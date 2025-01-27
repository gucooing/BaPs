package game

import (
	"strconv"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetFriendBin(s *enter.Session) *enter.AccountFriend {
	if s == nil {
		return nil
	}
	return s.AccountFriend
}

func GetFriendNum(s *enter.Session) int64 {
	bin := GetFriendBin(s)
	if bin == nil {
		return 0
	}
	return int64(len(bin.FriendList))
}

func GetAutoAcceptFriendRequest(s *enter.Session) bool {
	bin := GetFriendBin(s)
	if bin == nil {
		return false
	}
	return bin.AutoAcceptFriendRequest
}

func SetAutoAcceptFriendRequest(s *enter.Session, is bool) bool {
	bin := GetFriendBin(s)
	if bin == nil {
		return false
	}
	bin.AutoAcceptFriendRequest = is
	return true
}

func AddFriendByUid(s *enter.Session, uid int64) {
	bin := GetFriendBin(s)
	if bin == nil {
		return
	}
	if bin.FriendList == nil {
		bin.FriendList = make(map[int64]bool)
	}
	if bin.ReceivedList == nil {
		bin.ReceivedList = make(map[int64]bool)
	}
	if bin.SendReceivedList == nil {
		bin.SendReceivedList = make(map[int64]bool)
	}
	if bin.BlockedList == nil {
		bin.BlockedList = make(map[int64]bool)
	}
	bin.FriendList[uid] = true
	delete(bin.ReceivedList, uid)
	delete(bin.SendReceivedList, uid)
	delete(bin.BlockedList, uid)
}

func RemoveFriendByUid(s *enter.Session, uid int64) {
	bin := GetFriendBin(s)
	if bin == nil {
		return
	}
	if bin.FriendList == nil {
		bin.FriendList = make(map[int64]bool)
	}
	if bin.ReceivedList == nil {
		bin.ReceivedList = make(map[int64]bool)
	}
	if bin.SendReceivedList == nil {
		bin.SendReceivedList = make(map[int64]bool)
	}
	if bin.BlockedList == nil {
		bin.BlockedList = make(map[int64]bool)
	}
	delete(bin.FriendList, uid)
	delete(bin.ReceivedList, uid)
	delete(bin.SendReceivedList, uid)
}

func GetFriendDB(s *enter.Session) *proto.FriendDB {
	bin := GetBaseBin(s)
	if bin == nil {
		return nil
	}
	info := &proto.FriendDB{
		AccountId:                   bin.AccountId,
		Level:                       GetAccountLevel(s),
		Nickname:                    GetNickname(s),
		LastConnectTime:             mx.Unix(bin.LastConnectTime, 0),
		RepresentCharacterUniqueId:  GetRepresentCharacterUniqueId(s),
		RepresentCharacterCostumeId: 0,
		ComfortValue:                0,
		FriendCount:                 GetFriendNum(s),
		AttachmentDB:                GetAccountAttachmentDB(s),
	}

	return info
}

// GetFriendDBs 会进行玩家线程加锁操作,使用时需注意
func GetFriendDBs(s *enter.Session, uidList map[int64]bool) []*proto.FriendDB {
	if s == nil {
		return nil
	}
	list := make([]*proto.FriendDB, 0)
	for uid, ok := range uidList {
		// 跳过申请者,避免锁死
		if !ok || uid == s.AccountServerId {
			continue
		}
		friendS := enter.GetSessionByUid(uid)
		if friendS == nil {
			continue
		}
		friendS.GoroutinesSync.Lock()
		list = append(list, GetFriendDB(friendS))
		friendS.GoroutinesSync.Unlock()
	}

	return list
}

func GetFriendIdCardDB(s *enter.Session) *proto.FriendIdCardDB {
	bin := GetBaseBin(s)
	if bin == nil {
		return nil
	}
	info := &proto.FriendIdCardDB{
		Level:                       GetAccountLevel(s),
		Comment:                     GetComment(s),
		RepresentCharacterUniqueId:  GetRepresentCharacterUniqueId(s),
		LastConnectTime:             mx.Unix(bin.LastConnectTime, 0),
		FriendCode:                  strconv.FormatInt(s.AccountServerId, 10),
		EmblemId:                    GetEmblemUniqueId(s),
		CardBackgroundId:            GetCardBackgroundId(s),
		AutoAcceptFriendRequest:     GetAutoAcceptFriendRequest(s),
		SearchPermission:            GetSearchPermission(s),
		ShowAccountLevel:            GetShowAccountLevel(s),
		ShowArenaRanking:            GetShowArenaRanking(s),
		ArenaRanking:                1,
		ShowEliminateRaidRanking:    GetShowEliminateRaidRanking(s),
		EliminateRaidRanking:        2, // 前大决战排名
		EliminateRaidTier:           3,
		ShowFriendCode:              GetShowFriendCode(s),
		ShowRaidRanking:             GetShowRaidRanking(s),
		RaidRanking:                 4, // 前总力战排名
		RaidTier:                    5,
		RepresentCharacterCostumeId: 0,
	}

	return info
}

func GetDetailedAccountInfoDB(s *enter.Session) *proto.DetailedAccountInfoDB {
	if s == nil {
		return nil
	}
	info := &proto.DetailedAccountInfoDB{
		AccountId:                      s.AccountServerId,
		Nickname:                       GetNickname(s),
		Level:                          int64(GetAccountLevel(s)),
		Comment:                        GetComment(s),
		FriendCount:                    GetFriendNum(s),
		FriendCode:                     strconv.FormatInt(s.AccountServerId, 10),
		RepresentCharacterUniqueId:     GetRepresentCharacterUniqueId(s),
		ClanName:                       "",
		CharacterCount:                 GetCharacterCount(s), // 学生数量
		LastNormalCampaignClearStageId: 0,
		LastHardCampaignClearStageId:   0,
		ArenaRanking:                   1,
		RaidRanking:                    4,
		RaidTier:                       3,
		EliminateRaidRanking:           2,
		EliminateRaidTier:              3,
		AssistCharacterDBs:             make([]*proto.AssistCharacterDB, 0),
	}
	return info
}
