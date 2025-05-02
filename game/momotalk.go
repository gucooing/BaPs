package game

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewMomoTalkBin() *sro.MomoTalkBin {
	return &sro.MomoTalkBin{}
}

func GetMomoTalkBin(s *enter.Session) *sro.MomoTalkBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.MomoTalkBin == nil {
		bin.MomoTalkBin = NewMomoTalkBin()
	}
	return bin.MomoTalkBin
}

func GetFavorScheduleInfoList(s *enter.Session) map[int64]*sro.FavorScheduleInfo {
	bin := GetMomoTalkBin(s)
	if bin == nil {
		return nil
	}
	if bin.FavorScheduleInfoList == nil {
		bin.FavorScheduleInfoList = make(map[int64]*sro.FavorScheduleInfo)
	}
	return bin.FavorScheduleInfoList
}

func GetFavorScheduleInfo(s *enter.Session, characterId int64) *sro.FavorScheduleInfo {
	bin := GetFavorScheduleInfoList(s)
	if bin[characterId] == nil {
		bin[characterId] = &sro.FavorScheduleInfo{
			ScheduleGroupList: make([]int64, 0),
			MomoTalkInfoList:  make(map[int64]*sro.MomoTalkInfo),
			CurMessageGroupId: 0,
		}
	}
	return bin[characterId]
}

func GetMemoryLobbyInfoList(s *enter.Session) map[int64]*sro.MemoryLobbyInfo {
	bin := GetMomoTalkBin(s)
	if bin == nil {
		return nil
	}
	if bin.MemoryLobbyInfoList == nil {
		bin.MemoryLobbyInfoList = make(map[int64]*sro.MemoryLobbyInfo)
	}
	return bin.MemoryLobbyInfoList
}

func GetMemoryLobbyInfo(s *enter.Session, memoryLobbyId int64) *sro.MemoryLobbyInfo {
	bin := GetMemoryLobbyInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[memoryLobbyId]
}

func UpMomoTalkInfo(s *enter.Session, characterId, gid, messageId int64) {
	// 验证是否有这个id
	conf := gdconf.GetAcademyMessangerExcelTable(gid)
	if conf == nil {
		return
	}
	// 简单验证是否允许
	if conf.MessageCondition == "FavorRankUp" {
		characterInfo := GetCharacterInfo(s, characterId)
		if characterInfo == nil ||
			int64(characterInfo.FavorRank) < conf.ConditionValue {
			return
		}
	}
	bin := GetFavorScheduleInfo(s, characterId)
	if bin == nil {
		return
	}
	if bin.MomoTalkInfoList == nil {
		bin.MomoTalkInfoList = make(map[int64]*sro.MomoTalkInfo)
	}
	bin.CurMessageGroupId = gid
	if info, ok := bin.MomoTalkInfoList[gid]; ok {
		info.MessageGroupId = alg.MaxInt64(info.MessageGroupId, gid)
		info.ChosenMessageId = alg.MaxInt64(info.ChosenMessageId, messageId)
		info.ChosenDate = time.Now().Unix()
	}
	bin.MomoTalkInfoList[gid] = &sro.MomoTalkInfo{
		MessageGroupId:  gid,
		ChosenMessageId: messageId,
		ChosenDate:      time.Now().Unix(),
	}
}

func UpScheduleGroup(s *enter.Session, scheduleId int64) []*ParcelResult {
	conf := gdconf.GetAcademyFavorScheduleExcelTable(scheduleId)
	if conf == nil || s == nil {
		return nil
	}
	bin := GetFavorScheduleInfo(s, conf.CharacterId)
	if bin == nil {
		return nil
	}
	if bin.ScheduleGroupList == nil {
		bin.ScheduleGroupList = make([]int64, 0)
	}
	for _, id := range bin.ScheduleGroupList {
		if id == scheduleId {
			return nil
		}
	}
	bin.ScheduleGroupList = append(bin.ScheduleGroupList, scheduleId)

	return GetParcelResultList(conf.RewardParcelType, conf.RewardParcelId, conf.RewardAmount, false)
}

func UpMemoryLobbyInfo(s *enter.Session, memoryLobbyUniqueId int64) {
	bin := GetMomoTalkBin(s)
	if bin == nil {
		return
	}
	if bin.MemoryLobbyInfoList == nil {
		bin.MemoryLobbyInfoList = make(map[int64]*sro.MemoryLobbyInfo)
	}
	bin.MemoryLobbyInfoList[memoryLobbyUniqueId] = &sro.MemoryLobbyInfo{
		MemoryLobbyId: memoryLobbyUniqueId,
		ChosenDate:    time.Now().Unix(),
	}
}

func GetMomoTalkChoiceDBs(s *enter.Session, characterId int64) []*proto.MomoTalkChoiceDB {
	list := make([]*proto.MomoTalkChoiceDB, 0)
	characterInfo := GetCharacterInfo(s, characterId)
	bin := GetFavorScheduleInfo(s, characterId)
	if bin == nil || characterInfo == nil {
		return list
	}
	for _, info := range bin.MomoTalkInfoList {
		list = append(list, &proto.MomoTalkChoiceDB{
			CharacterDBId:   characterInfo.ServerId,
			MessageGroupId:  info.MessageGroupId,
			ChosenMessageId: info.ChosenMessageId,
			ChosenDate:      mx.Unix(info.ChosenDate, 0),
		})
	}
	return list
}

func GetMemoryLobbyDBs(s *enter.Session) []*proto.MemoryLobbyDB {
	list := make([]*proto.MemoryLobbyDB, 0)
	if s == nil {
		return list
	}
	for _, info := range GetMemoryLobbyInfoList(s) {
		list = append(list, &proto.MemoryLobbyDB{
			Type:                proto.ParcelType_MemoryLobby,
			MemoryLobbyUniqueId: info.MemoryLobbyId,
		})
	}

	return list
}

func GetMemoryLobbyDB(s *enter.Session, memoryLobbyId int64) *proto.MemoryLobbyDB {
	bin := GetMemoryLobbyInfo(s, memoryLobbyId)
	if bin == nil {
		return nil
	}
	return &proto.MemoryLobbyDB{
		Type:                proto.ParcelType_MemoryLobby,
		MemoryLobbyUniqueId: bin.MemoryLobbyId,
	}
}
