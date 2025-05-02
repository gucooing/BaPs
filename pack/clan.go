package pack

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func ClanLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ClanLoginResponse)

	rsp.AccountClanDB = game.GetClanDB(enter.GetYostarClanByServerId(game.GetClanServerId(s))) // 社团简介
	rsp.AccountClanMemberDB = game.GetClanMemberDB(s)                                          // 本人信息
	rsp.ClanAssistSlotDBs = make([]*proto.ClanAssistSlotDB, 0)                                 // 援助信息
}

func ClanCheck(s *enter.Session, request, response mx.Message) {

}

func ClanLobby(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ClanLobbyResponse)

	game.SetLastLoginTime(s)

	clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
	rsp.IrcConfig = game.GetIrcServerConfig(s)
	rsp.AccountClanMemberDB = game.GetClanMemberDB(s) // 本人信息
	rsp.AccountClanDB = game.GetClanDB(clanInfo)
	rsp.ClanMemberDBs = game.GetClanMemberDBs(clanInfo)
	rsp.DefaultExposedClanDBs = game.GetDefaultExposedClanDBs(s)
}

func ClanSearch(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanSearchRequest)
	rsp := response.(*proto.ClanSearchResponse)

	rsp.ClanDBs = make([]*proto.ClanDB, 0)
	if clanInfo := enter.GetYostarClanByClanName(req.SearchString); clanInfo != nil {
		rsp.ClanDBs = append(rsp.ClanDBs, game.GetClanDB(clanInfo))
	}
	if clanInfo := enter.GetYostarClanByServerId(alg.S2I64(req.ClanUniqueCode)); clanInfo != nil {
		rsp.ClanDBs = append(rsp.ClanDBs, game.GetClanDB(clanInfo))
	}
	if len(rsp.ClanDBs) == 0 {
		rsp.ClanDBs = game.GetDefaultExposedClanDBs(s)
	}
}

func ClanCreate(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanCreateRequest)
	rsp := response.(*proto.ClanCreateResponse)

	defer func() {
		rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
		rsp.ClanMemberDB = game.GetClanMemberDB(s)
		rsp.ClanDB = game.GetClanDB(enter.GetYostarClanByServerId(game.GetClanServerId(s)))
	}()

	s.Error = game.NewClan(s, req.ClanNickName, req.ClanJoinOption)
	if s.Error == 0 {
		game.UpCurrency(s, proto.CurrencyTypes_GemBonus, -100)
	}
}

func ClanMemberList(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanMemberListRequest)
	rsp := response.(*proto.ClanMemberListResponse)

	clanInfo := enter.GetYostarClanByServerId(req.ClanDBId)
	if clanInfo == nil {
		s.Error = 0
		return
	}

	rsp.ClanDB = game.GetClanDB(clanInfo)
	rsp.ClanMemberDBs = game.GetClanMemberDBs(clanInfo)
}

func ClanJoin(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanJoinRequest)
	rsp := response.(*proto.ClanJoinResponse)

	defer func() {
		rsp.IrcConfig = game.GetIrcServerConfig(s)
		rsp.ClanMemberDB = game.GetClanMemberDB(s)
	}()

	if enter.GetYostarClanByServerId(game.GetClanServerId(s)) != nil {
		// 已有社团
		return
	}
	clanInfo := enter.GetYostarClanByServerId(req.ClanDBId)
	if clanInfo == nil {
		return
	}
	rsp.ClanDB = game.GetClanDB(clanInfo)
	if clanInfo.GetMemberCount() >= enter.ClanMaxMemberCount {
		// 满人了
		return
	}
	if clanInfo.JoinOption == int32(proto.ClanJoinOption_Free) {
		// 自动加入
		clanInfo.AddAccount(s.AccountServerId, int32(proto.ClanSocialGrade_Member))
		game.SetClanServerId(s, clanInfo.ServerId)
		rsp.ClanMemberDB = game.GetClanMemberDB(s) // 本人信息
		return
	}
	clanInfo.AddApplicantAccount(s.AccountServerId)
}

func ClanAutoJoin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ClanJoinResponse)

	defer func() {
		clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
		rsp.IrcConfig = game.GetIrcServerConfig(s)
		rsp.ClanMemberDB = game.GetClanMemberDB(s) // 本人信息
		rsp.ClanDB = game.GetClanDB(clanInfo)
	}()
	// 自动加入逻辑实现
	for _, clanInfo := range enter.GetAllYostarClanList() {
		if clanInfo.JoinOption == int32(proto.ClanJoinOption_Free) &&
			clanInfo.GetMemberCount() < int64(enter.ClanMaxMemberCount) {
			clanInfo.AddAccount(s.AccountServerId, int32(proto.ClanSocialGrade_Member))
			break
		}
	}
}

func ClanSetting(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanSettingRequest)
	rsp := response.(*proto.ClanSettingResponse)

	clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
	defer func() {
		rsp.ClanDB = game.GetClanDB(clanInfo)
	}()

	clanInfo.SetNotice(req.ChangedNotice)
	clanInfo.SetJoinOption(int32(req.ClanJoinOption))
	// name
	if clanInfo.ClanName != req.ChangedClanName {
		oldName := clanInfo.ClanName
		clanInfo.ClanName = req.ChangedClanName
		if clanInfo.UpDate() != nil {
			logger.Debug("社团名称修改失败,应该是有重复的名称,ClanName:%s", clanInfo.ClanName)
			clanInfo.ClanName = oldName
		}
	}
}

func ClanApplicant(s *enter.Session, request, response mx.Message) {
	// req := request.(*proto.ClanApplicantRequest)
	rsp := response.(*proto.ClanApplicantResponse)

	rsp.ClanMemberDBs = make([]*proto.ClanMemberDB, 0)
	clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
	for _, ca := range clanInfo.GetAllApplicantAccount() {
		ps := enter.GetSessionByUid(ca.Uid)
		rsp.ClanMemberDBs = append(rsp.ClanMemberDBs, &proto.ClanMemberDB{
			AccountId:                   ps.AccountServerId,
			AccountLevel:                int64(game.GetAccountLevel(ps)),
			AccountNickName:             game.GetNickname(ps),
			AttachmentDB:                game.GetAccountAttachmentDB(ps),
			CafeComfortValue:            9000,
			RepresentCharacterUniqueId:  game.GetRepresentCharacterUniqueId(ps),
			GameLoginDate:               game.GetLastConnectTime(ps),
			ClanDBId:                    clanInfo.ServerId,
			RepresentCharacterCostumeId: 0,
			ClanSocialGrade:             proto.ClanSocialGrade_Applicant,
			AppliedDate:                 mx.Unix(ca.ApplicantTime, 0),
		})
	}
}

func ClanMember(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanMemberRequest)
	rsp := response.(*proto.ClanMemberResponse)

	clanInfo := enter.GetYostarClanByServerId(req.ClanDBId)
	rsp.ClanDB = game.GetClanDB(clanInfo)

	ps := enter.GetSessionByUid(req.MemberAccountId)
	rsp.ClanMemberDB = game.GetClanMemberDB(ps)
	rsp.DetailedAccountInfoDB = game.GetDetailedAccountInfoDB(ps, proto.AssistRelation_Clan)
}

func ClanQuit(s *enter.Session, request, response mx.Message) {
	clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
	if clanInfo == nil {
		return
	}
	clanInfo.RemoveAccount(s.AccountServerId)
	game.SetClanServerId(s, 0)
}

func ClanKick(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanKickRequest)

	clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
	if clanInfo == nil {
		return
	}
	myAc := clanInfo.GetClanAccount(s.AccountServerId)
	if myAc == nil ||
		(myAc.SocialGrade != int32(proto.ClanSocialGrade_President) &&
			myAc.SocialGrade != int32(proto.ClanSocialGrade_Manager)) {
		return
	}
	ps := enter.GetSessionByUid(req.MemberAccountId)
	clanInfo.RemoveAccount(req.MemberAccountId)
	game.SetClanServerId(ps, 0)
}

func ClanConfer(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanConferRequest)
	rsp := response.(*proto.ClanConferResponse)

	clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
	defer func() {
		rsp.ClanDB = game.GetClanDB(clanInfo)
		rsp.ClanMemberDB = game.GetClanMemberDB(s)
	}()
	myAc := clanInfo.GetClanAccount(s.AccountServerId)
	fs := enter.GetSessionByUid(req.MemberAccountId)
	if myAc == nil || fs == nil ||
		myAc.SocialGrade > int32(req.ConferingGrade) {
		return
	}
	rsp.AccountClanMemberDB = game.GetClanMemberDB(fs)
	if clanInfo.GetClanAccount(req.MemberAccountId) == nil {
		return
	}
	clanInfo.SetPresident(req.MemberAccountId)
}

func ClanPermit(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanPermitRequest)
	rsp := response.(*proto.ClanPermitResponse)

	clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
	defer func() {
		rsp.ClanDB = game.GetClanDB(clanInfo)
	}()
	myAc := clanInfo.GetClanAccount(s.AccountServerId)
	if myAc == nil || myAc.SocialGrade > int32(proto.ClanSocialGrade_Manager) ||
		myAc.SocialGrade == int32(proto.ClanSocialGrade_None) {
		return
	}
	aacList := clanInfo.GetAllApplicantAccount()
	fs := enter.GetSessionByUid(req.ApplicantAccountId)
	if fs == nil || aacList[req.ApplicantAccountId] == nil ||
		enter.GetYostarClanByServerId(game.GetClanServerId(fs)) != nil {
		return
	}
	clanInfo.RemoveApplicantAccount(req.ApplicantAccountId)
	if req.IsPerMit {
		clanInfo.AddAccount(req.ApplicantAccountId, int32(proto.ClanSocialGrade_Member))
		game.SetClanServerId(fs, clanInfo.ServerId)
	}
	rsp.ClanMemberDB = game.GetClanMemberDB(fs)
}

func ClanMyAssistList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ClanMyAssistListResponse)

	rsp.ClanAssistSlotDBs = game.GetClanAssistSlotDBs(s)
}

func ClanSetAssist(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanSetAssistRequest)
	rsp := response.(*proto.ClanSetAssistResponse)

	rewardInfo := &proto.ClanAssistRewardInfo{
		CharacterDBId:            0,
		DeployDate:               mx.MxTime{},
		CumultativeRewardParcels: make([]*proto.ParcelInfo, 0),
	}
	rsp.RewardInfo = rewardInfo
	parcelResultList := make([]*game.ParcelResult, 0)
	defer func() {
		rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
	}()

	bin := game.GetAssistList(s)
	characterInfo := s.GetCharacterByKeyId(req.CharacterDBId)
	if bin == nil || characterInfo == nil {
		return
	}
	if bin[int32(req.EchelonType)] == nil {
		bin[int32(req.EchelonType)] = &sro.AssistList{
			AssistInfoList: make(map[int32]*sro.AssistInfo),
		}
	}
	assistList := bin[int32(req.EchelonType)]
	if assistList.AssistInfoList == nil {
		assistList.AssistInfoList = make(map[int32]*sro.AssistInfo)
	}
	info := &sro.AssistInfo{
		EchelonType:    int32(req.EchelonType),
		SlotNumber:     int64(req.SlotNumber),
		CharacterId:    characterInfo.CharacterId,
		DeployDate:     time.Now().Unix(),
		TotalRentCount: 0,
	}
	old := assistList.AssistInfoList[req.SlotNumber]
	assistList.AssistInfoList[req.SlotNumber] = info
	rsp.ClanAssistSlotDB = game.GetClanAssistSlotDB(s, info)
	if old == nil {
		return
	}
	rewardInfo.CharacterDBId = game.GetCharacterInfo(s, old.CharacterId).GetServerId()
	rewardInfo.DeployDate = mx.Unix(old.DeployDate, 0)
	num := time.Now().Sub(time.Unix(old.DeployDate, 0)) / time.Minute
	rentCount := alg.MinInt64(int64(num*game.AssistTermRewardPeriodFromSec)+
		alg.MinInt64(old.TotalRentCount, game.AssistRentRewardDailyMaxCount)*game.AssistRentalFeeAmount,
		game.AssistRewardLimit)
	rewardInfo.CumultativeRewardParcels = append(rewardInfo.CumultativeRewardParcels,
		game.GetParcelInfo(int64(proto.CurrencyTypes_Gold), rentCount, proto.ParcelType_Currency))
	parcelResultList = append(parcelResultList, &game.ParcelResult{
		ParcelType: proto.ParcelType_Currency,
		ParcelId:   int64(proto.CurrencyTypes_Gold),
		Amount:     rentCount,
	})
}

func ClanAllAssistList(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ClanAllAssistListRequest)
	rsp := response.(*proto.ClanAllAssistListResponse)

	rsp.AssistCharacterDBs = make([]*proto.AssistCharacterDB, 0)
	rsp.AssistCharacterRentHistoryDBs = make([]*proto.ClanAssistRentHistoryDB, 0)
	rsp.ClanDBId = game.GetClanServerId(s)

	echelonType := req.EchelonType
	switch echelonType {
	case proto.EchelonType_EliminateRaid01,
		proto.EchelonType_EliminateRaid02,
		proto.EchelonType_EliminateRaid03,
		proto.EchelonType_MultiFloorRaid:
		echelonType = proto.EchelonType_Raid
	}

	addAssistCharacter := func(fs *enter.Session, assistRelation proto.AssistRelation) {
		assist := game.GetAssistListByEchelonType(fs, echelonType)
		for _, info := range assist.GetAssistInfoList() {
			if info == nil {
				continue
			}
			characterInfo := game.GetCharacterInfo(fs, info.CharacterId)
			if characterInfo == nil {
				continue
			}
			rsp.AssistCharacterDBs = append(rsp.AssistCharacterDBs,
				game.GetAssistCharacterDB(fs, info, assistRelation))

			// 已借的人历史数据
			// rsp.AssistCharacterRentHistoryDBs = append(rsp.AssistCharacterRentHistoryDBs,
			// 	&proto.ClanAssistRentHistoryDB{
			// 		AssistCharacterAccountId: fs.AccountServerId,
			// 		AssistCharacterDBId:      characterInfo.ServerId,
			// 		// RentDate:                 mx.Unix(info.DeployDate, 0), 是我的视角下什么时候借了人
			// 		AssistCharacterId: characterInfo.CharacterId,
			// 	})
		}
	}
	addUid := make(map[int64]bool, 0)
	// 添加好友的
	for uid, _ := range game.GetFriendBin(s).GetFriendList() {
		if uid == s.AccountServerId || addUid[uid] {
			continue
		}
		fs := enter.GetSessionByUid(uid)
		if fs == nil {
			continue
		}
		addAssistCharacter(fs, proto.AssistRelation_Friend)
		addUid[uid] = true
	}
	// 添加社团
	clanInfo := enter.GetYostarClanByServerId(game.GetClanServerId(s))
	for uid, _ := range clanInfo.GetAllAccount() {
		if uid == s.AccountServerId || addUid[uid] {
			continue
		}
		fs := enter.GetSessionByUid(uid)
		if fs == nil {
			continue
		}
		addAssistCharacter(fs, proto.AssistRelation_Clan)
		addUid[uid] = true
	}
}
