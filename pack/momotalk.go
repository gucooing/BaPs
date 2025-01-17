package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func MemoryLobbyList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MemoryLobbyListResponse)

	rsp.MemoryLobbyDBs = game.GetMemoryLobbyDBs(s)
}

func MomoTalkOutLine(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MomoTalkOutLineResponse)

	rsp.MomoTalkOutLineDBs = make([]*proto.MomoTalkOutLineDB, 0)
	rsp.FavorScheduleRecords = make(map[int64][]int64)

	for characterId, fsInfo := range game.GetFavorScheduleInfoList(s) {
		characterInfo := game.GetCharacterInfo(s, characterId)
		if characterInfo == nil {
			continue
		}
		if rsp.FavorScheduleRecords[characterId] == nil {
			rsp.FavorScheduleRecords[characterId] = make([]int64, 0)
		}
		// 添加已完成的
		for gid, ok := range fsInfo.ScheduleGroupList {
			if ok {
				rsp.FavorScheduleRecords[characterId] = append(
					rsp.FavorScheduleRecords[characterId], gid)
			}
		}
		// 添加最新的
		if info, ok := fsInfo.MomoTalkInfoList[fsInfo.CurMessageGroupId]; ok {
			rsp.MomoTalkOutLineDBs = append(rsp.MomoTalkOutLineDBs, &proto.MomoTalkOutLineDB{
				CharacterDBId:        characterInfo.ServerId,
				CharacterId:          characterId,
				LatestMessageGroupId: info.MessageGroupId,
				ChosenMessageId:      info.ChosenMessageId,
				LastUpdateDate:       mx.Unix(info.ChosenDate, 0),
			})
		}
	}
}

func MomoTalkMessageList(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.MomoTalkMessageListRequest)
	rsp := response.(*proto.MomoTalkMessageListResponse)

	characterId := game.ServerIdToCharacterId(s, req.CharacterDBId)

	rsp.MomoTalkOutLineDB = &proto.MomoTalkOutLineDB{
		CharacterDBId:        req.CharacterDBId,
		CharacterId:          characterId,
		LatestMessageGroupId: 0,
		ChosenMessageId:      0,
		LastUpdateDate:       mx.MxTime{},
	}
	rsp.MomoTalkChoiceDBs = game.GetMomoTalkChoiceDBs(s, characterId)

	bin := game.GetFavorScheduleInfo(s, characterId)
	if bin == nil {
		return
	}
	if info, ok := bin.MomoTalkInfoList[bin.CurMessageGroupId]; ok {
		rsp.MomoTalkOutLineDB.LatestMessageGroupId = info.MessageGroupId
		rsp.MomoTalkOutLineDB.ChosenMessageId = info.ChosenMessageId
		rsp.MomoTalkOutLineDB.LastUpdateDate = mx.Unix(info.ChosenDate, 0)
	}
}

func MomoTalkRead(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.MomoTalkReadRequest)
	rsp := response.(*proto.MomoTalkReadResponse)

	characterId := game.ServerIdToCharacterId(s, req.CharacterDBId)

	game.UpMomoTalkInfo(s, characterId, req.LastReadMessageGroupId, req.ChosenMessageId)

	rsp.MomoTalkChoiceDBs = game.GetMomoTalkChoiceDBs(s, characterId)
	rsp.MomoTalkOutLineDB = &proto.MomoTalkOutLineDB{
		CharacterDBId:        req.CharacterDBId,
		CharacterId:          characterId,
		LatestMessageGroupId: req.LastReadMessageGroupId,
		ChosenMessageId:      req.ChosenMessageId,
		LastUpdateDate:       mx.Now(),
	}
}

func MomoTalkFavorSchedule(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.MomoTalkFavorScheduleRequest)
	rsp := response.(*proto.MomoTalkFavorScheduleResponse)

	parcelResultList := game.UpScheduleGroup(s, req.ScheduleId)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
	rsp.FavorScheduleRecords = make(map[int64][]int64)
	for characterId, fsInfo := range game.GetFavorScheduleInfoList(s) {
		if rsp.FavorScheduleRecords[characterId] == nil {
			rsp.FavorScheduleRecords[characterId] = make([]int64, 0)
		}
		for gid, ok := range fsInfo.ScheduleGroupList {
			if ok {
				rsp.FavorScheduleRecords[characterId] = append(
					rsp.FavorScheduleRecords[characterId], gid)
			}
		}
	}
}
