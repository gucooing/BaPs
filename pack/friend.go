package pack

import (
	"math/rand"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func FriendCheck(s *enter.Session, request, response proto.Message) {
	af := s.AccountFriend
	if af == nil {
		logger.Warn("[UID:%v]好友数据拉取失败,请查询数据库是否正常", s.AccountServerId)
		return
	}
	if len(af.ReceivedList) > 0 {
		game.SetServerNotification(s, proto.ServerNotificationFlag_HasFriendRequest, true)
	} else {
		game.SetServerNotification(s, proto.ServerNotificationFlag_HasFriendRequest, false)
	}
}

func FriendList(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.FriendListResponse)

	bin := game.GetFriendBin(s)
	if bin == nil {
		return
	}
	bin.SyncAf.RLock()
	defer bin.SyncAf.RUnlock()

	rsp.IdCardBackgroundDBs = game.GetIdCardBackgroundDBs(s)
	rsp.FriendDBs = game.GetFriendDBs(s, bin.FriendList)
	rsp.SentRequestFriendDBs = game.GetFriendDBs(s, bin.SendReceivedList)
	rsp.ReceivedRequestFriendDBs = game.GetFriendDBs(s, bin.ReceivedList)
	rsp.BlockedUserDBs = game.GetFriendDBs(s, bin.BlockedList)
	rsp.FriendIdCardDB = game.GetFriendIdCardDB(s)
}

func FriendGetIdCard(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.FriendGetIdCardResponse)

	rsp.FriendIdCardDB = game.GetFriendIdCardDB(s)
}

func FriendSetIdCard(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.FriendSetIdCardRequest)

	game.SetComment(s, req.Comment)
	game.SetCardRepresentCharacterUniqueId(s, req.RepresentCharacterUniqueId)
	game.SetEmblemUniqueId(s, req.EmblemId)
	game.SetCardBackgroundId(s, req.BackgroundId)
	game.SetAutoAcceptFriendRequest(s, req.AutoAcceptFriendRequest)
	game.SetSearchPermission(s, req.SearchPermission)
	game.SetShowAccountLevel(s, req.ShowAccountLevel)
	game.SetShowFriendCode(s, req.ShowFriendCode)
	game.SetShowRaidRanking(s, req.ShowRaidRanking)
	game.SetShowArenaRanking(s, req.ShowArenaRanking)
	game.SetShowEliminateRaidRanking(s, req.ShowEliminateRaidRanking)
}

func FriendSearch(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.FriendSearchRequest)
	rsp := response.(*proto.FriendSearchResponse)

	rsp.SearchResult = make([]*proto.FriendDB, 0)
	// 搜索玩家
	if uid := alg.S2I64(req.FriendCode); uid != 0 && uid != s.AccountServerId {
		friendS := enter.GetSessionByUid(uid)
		rsp.SearchResult = append(rsp.SearchResult, game.GetFriendDB(friendS))
		return
	}

	allSession := enter.GetAllSessionList()
	maxNum := alg.MainInt(30, len(allSession)-1)
	uidlist := make(map[int64]bool, 0)
	for i := 0; i < maxNum; i++ {
		friendS := allSession[rand.Intn(len(allSession))]
		if friendS == nil ||
			friendS.AccountServerId == s.AccountServerId {
			continue
		}
		uidlist[friendS.AccountServerId] = true
	}
	rsp.SearchResult = game.GetFriendDBs(s, uidlist)
}

func FriendGetFriendDetailedInfo(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.FriendGetFriendDetailedInfoRequest)
	rsp := response.(*proto.FriendGetFriendDetailedInfoResponse)

	friendS := enter.GetSessionByUid(req.FriendAccountId)
	if friendS == nil {
		return
	}
	rsp.AttachmentDB = game.GetAccountAttachmentDB(friendS)
	rsp.DetailedAccountInfoDB = game.GetDetailedAccountInfoDB(friendS, proto.AssistRelation_Friend)
}

func FriendSendFriendRequest(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.FriendSendFriendRequestRequest)
	rsp := response.(*proto.FriendSendFriendRequestResponse)

	bin := game.GetFriendBin(s)
	if bin == nil {
		return
	}
	bin.SyncAf.Lock()
	defer bin.SyncAf.Unlock()

	defer func() {
		rsp.FriendDBs = game.GetFriendDBs(s, bin.FriendList)
		rsp.SentRequestFriendDBs = game.GetFriendDBs(s, bin.SendReceivedList)
		rsp.ReceivedRequestFriendDBs = game.GetFriendDBs(s, bin.ReceivedList)
		rsp.BlockedUserDBs = game.GetFriendDBs(s, bin.BlockedList)
	}()
	if _, ok := bin.FriendList[req.TargetAccountId]; ok {
		return
	}
	friendS := enter.GetSessionByUid(req.TargetAccountId)
	if friendS == nil {
		return
	}
	targetFriendBin := game.GetFriendBin(friendS)
	if targetFriendBin == nil {
		logger.Warn("[UID:%v]好友信息拉取失败,请检查数据库连接情况", req.TargetAccountId)
		return
	}
	targetFriendBin.SyncAf.Lock()
	defer targetFriendBin.SyncAf.Unlock()
	// 直接成为好友
	if _, ok := targetFriendBin.SendReceivedList[s.AccountServerId]; ok ||
		targetFriendBin.AutoAcceptFriendRequest {
		game.AddFriendByUid(friendS, s.AccountServerId)
		game.AddFriendByUid(s, req.TargetAccountId)
		return
	}
	// 不是则添加到待录取列表中
	if bin.SendReceivedList == nil {
		bin.SendReceivedList = make(map[int64]bool)
	}
	bin.SendReceivedList[req.TargetAccountId] = true
	if targetFriendBin.ReceivedList == nil {
		targetFriendBin.ReceivedList = make(map[int64]bool)
	}
	targetFriendBin.ReceivedList[s.AccountServerId] = true
}

func FriendAcceptFriendRequest(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.FriendAcceptFriendRequestRequest)
	rsp := response.(*proto.FriendAcceptFriendRequestResponse)

	bin := game.GetFriendBin(s)
	if bin == nil {
		return
	}
	bin.SyncAf.Lock()
	defer bin.SyncAf.Unlock()

	defer func() {
		rsp.FriendDBs = game.GetFriendDBs(s, bin.FriendList)
		rsp.SentRequestFriendDBs = game.GetFriendDBs(s, bin.SendReceivedList)
		rsp.ReceivedRequestFriendDBs = game.GetFriendDBs(s, bin.ReceivedList)
		rsp.BlockedUserDBs = game.GetFriendDBs(s, bin.BlockedList)
	}()

	friendS := enter.GetSessionByUid(req.TargetAccountId)
	if friendS == nil {
		return
	}
	targetFriendBin := game.GetFriendBin(friendS)
	if targetFriendBin == nil {
		logger.Warn("[UID:%v]好友信息拉取失败,请检查数据库连接情况", req.TargetAccountId)
		return
	}
	targetFriendBin.SyncAf.Lock()
	defer targetFriendBin.SyncAf.Unlock()
	game.AddFriendByUid(friendS, s.AccountServerId)
	game.AddFriendByUid(s, req.TargetAccountId)
}

func FriendDeclineFriendRequest(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.FriendDeclineFriendRequestRequest)
	rsp := response.(*proto.FriendDeclineFriendRequestResponse)

	bin := game.GetFriendBin(s)
	if bin == nil {
		return
	}
	bin.SyncAf.Lock()
	defer bin.SyncAf.Unlock()

	defer func() {
		rsp.FriendDBs = game.GetFriendDBs(s, bin.FriendList)
		rsp.SentRequestFriendDBs = game.GetFriendDBs(s, bin.SendReceivedList)
		rsp.ReceivedRequestFriendDBs = game.GetFriendDBs(s, bin.ReceivedList)
		rsp.BlockedUserDBs = game.GetFriendDBs(s, bin.BlockedList)
	}()
	friendS := enter.GetSessionByUid(req.TargetAccountId)
	if friendS == nil {
		return
	}
	targetFriendBin := game.GetFriendBin(friendS)
	if targetFriendBin == nil {
		logger.Warn("[UID:%v]好友信息拉取失败,请检查数据库连接情况", req.TargetAccountId)
		return
	}
	targetFriendBin.SyncAf.Lock()
	defer targetFriendBin.SyncAf.Unlock()

	if bin.SendReceivedList == nil {
		bin.SendReceivedList = make(map[int64]bool)
	}
	delete(bin.SendReceivedList, req.TargetAccountId)
	if targetFriendBin.ReceivedList == nil {
		targetFriendBin.ReceivedList = make(map[int64]bool)
	}
	delete(targetFriendBin.ReceivedList, s.AccountServerId)
}

func FriendRemove(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.FriendRemoveRequest)
	rsp := response.(*proto.FriendRemoveResponse)

	bin := game.GetFriendBin(s)
	if bin == nil {
		return
	}
	bin.SyncAf.Lock()
	defer bin.SyncAf.Unlock()

	defer func() {
		rsp.FriendDBs = game.GetFriendDBs(s, bin.FriendList)
		rsp.SentRequestFriendDBs = game.GetFriendDBs(s, bin.SendReceivedList)
		rsp.ReceivedRequestFriendDBs = game.GetFriendDBs(s, bin.ReceivedList)
		rsp.BlockedUserDBs = game.GetFriendDBs(s, bin.BlockedList)
	}()
	friendS := enter.GetSessionByUid(req.TargetAccountId)
	if friendS == nil {
		return
	}
	targetFriendBin := game.GetFriendBin(friendS)
	if targetFriendBin == nil {
		logger.Warn("[UID:%v]好友信息拉取失败,请检查数据库连接情况", req.TargetAccountId)
		return
	}
	targetFriendBin.SyncAf.Lock()
	defer targetFriendBin.SyncAf.Unlock()
	game.RemoveFriendByUid(s, req.TargetAccountId)
	game.RemoveFriendByUid(friendS, s.AccountServerId)
}
