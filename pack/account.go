package pack

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"strconv"
	"time"

	"github.com/gucooing/BaPs/common/check"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/cmd"
	"github.com/gucooing/BaPs/protocol/proto"
)

func AccountAuth(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountAuthRequest)
	rsp := response.(*proto.AccountAuthResponse)

	rsp.CurrentVersion = req.Version
	rsp.BattleValidation = true
	rsp.AccountDB = game.GetAccountDB(s)
	rsp.StaticOpenConditions = game.GetStaticOpenConditions(s)
	rsp.AttendanceBookRewards = game.GetAttendanceBookRewards(s)
	rsp.AttendanceHistoryDBs = game.GetAttendanceHistoryDBs(s)

	rsp.IssueAlertInfos = make([]*proto.IssueAlertInfoDB, 0)
	rsp.RepurchasableMonthlyProductCountDBs = make([]*proto.PurchaseCountDB, 0)
	rsp.MonthlyProductParcel = make([]*proto.ParcelInfo, 0)
	rsp.MonthlyProductMail = make([]*proto.ParcelInfo, 0)
	rsp.BiweeklyProductParcel = make([]*proto.ParcelInfo, 0)
	rsp.BiweeklyProductMail = make([]*proto.ParcelInfo, 0)
	rsp.WeeklyProductParcel = make([]*proto.ParcelInfo, 0)
	rsp.WeeklyProductMail = make([]*proto.ParcelInfo, 0)
	rsp.EncryptedUID = "1"
	rsp.AccountRestrictionsDB = &proto.AccountRestrictionsDB{}
	rsp.TTSCdnUri = "https://prod-voice.bluearchiveyostar.com/prod_new_2/version2/"

	s.AccountState = proto.AccountState_Normal
	game.SetLastConnectTime(s)

	game.AddToast(s, "欢迎游玩BaPs,这是一个半开源的免费服务器")
	// 任务二次处理
	mission := game.GetMissionBin(s)
	for t, info := range mission.GetCategoryMissionInfo() {
		if t == "" {
			delete(mission.GetCategoryMissionInfo(), t)
			continue
		}
		s.AddMissionByCompleteConditionType(info)
	}
}

func AccountNickname(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountNicknameRequest)
	rsp := response.(*proto.AccountNicknameResponse)

	if !check.CheckName(req.Nickname) {
		logger.Warn("玩家昵称检查不通过,但未作拦截处理,仅通知")
	}
	game.SetAccountNickname(s, req.Nickname)
	rsp.AccountDB = game.GetAccountDB(s)
}

func AccountCallName(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountCallNameRequest)
	rsp := response.(*proto.AccountCallNameResponse)

	game.SetCallName(s, req.CallName)
	rsp.AccountDB = game.GetAccountDB(s)
}

func ProofTokenRequestQuestion(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ProofTokenRequestQuestionResponse)

	rsp.Question = "1"
	rsp.Hint = 1
}

func NetworkTimeSync(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.NetworkTimeSyncResponse)

	rsp.ReceiveTick = game.GetServerTime()
	rsp.ServerTimeTicks = game.GetServerTime()
	rsp.EchoSendTick = game.GetServerTimeTick()
}

func AccountLoginSync(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountLoginSyncRequest)
	rsp := response.(*proto.AccountLoginSyncResponse)

	rsp.FriendCode = strconv.FormatInt(s.AccountServerId, 10)
	rsp.FriendCount = game.GetFriendNum(s)
	rsp.StaticOpenConditions = game.GetStaticOpenConditions(s)

	for _, cmdId := range req.SyncProtocols {
		syncReq := cmd.Get().GetRequestPacketByCmdId(cmdId)
		if syncReq == nil {
			logger.Error("AccountLoginSync SyncProtocol Req failed:%v", cmdId)
			continue
		}
		syncRsp := cmd.Get().GetResponsePacketByCmdId(cmdId)
		if syncRsp == nil {
			logger.Error("AccountLoginSync SyncProtocol Rsp failed:%v", cmdId)
			continue
		}
		// syncRsp.SetSessionKey(rsp.BasePacket)
		switch cmdId {
		case proto.Protocol_Cafe_Get:
			CafeGetInfo(s, syncReq, syncRsp)
			rsp.CafeGetInfoResponse = syncRsp.(*proto.CafeGetInfoResponse)
		case proto.Protocol_Account_CurrencySync:
			AccountCurrencySync(s, syncReq, syncRsp)
			rsp.AccountCurrencySyncResponse = syncRsp.(*proto.AccountCurrencySyncResponse)
		case proto.Protocol_Character_List:
			CharacterList(s, syncReq, syncRsp)
			rsp.CharacterListResponse = syncRsp.(*proto.CharacterListResponse)
		case proto.Protocol_Equipment_List:
			EquipmentList(s, syncReq, syncRsp)
			rsp.EquipmentItemListResponse = syncRsp.(*proto.EquipmentItemListResponse)
		case proto.Protocol_CharacterGear_List:
			CharacterGearList(s, syncReq, syncRsp)
			rsp.CharacterGearListResponse = syncRsp.(*proto.CharacterGearListResponse)
		case proto.Protocol_Echelon_List:
			EchelonList(s, syncReq, syncRsp)
			rsp.EchelonListResponse = syncRsp.(*proto.EchelonListResponse)
		case proto.Protocol_MemoryLobby_List:
			MemoryLobbyList(s, syncReq, syncRsp)
			rsp.MemoryLobbyListResponse = syncRsp.(*proto.MemoryLobbyListResponse)
		case proto.Protocol_Campaign_List:
			CampaignList(s, syncReq, syncRsp)
			rsp.CampaignListResponse = syncRsp.(*proto.CampaignListResponse)
		case proto.Protocol_Arena_Login:
			ArenaLogin(s, syncReq, syncRsp)
			rsp.ArenaLoginResponse = syncRsp.(*proto.ArenaLoginResponse)
		case proto.Protocol_Raid_Login:
			RaidLogin(s, syncReq, syncRsp)
			rsp.RaidLoginResponse = syncRsp.(*proto.RaidLoginResponse)
		case proto.Protocol_EliminateRaid_Login:
			EliminateRaidLogin(s, syncReq, syncRsp)
			rsp.EliminateRaidLoginResponse = syncRsp.(*proto.EliminateRaidLoginResponse)
		case proto.Protocol_Craft_List:
			CraftInfoList(s, syncReq, syncRsp)
			rsp.CraftInfoListResponse = syncRsp.(*proto.CraftInfoListResponse)
		case proto.Protocol_Clan_Login:
			ClanLogin(s, syncReq, syncRsp)
			rsp.ClanLoginResponse = syncRsp.(*proto.ClanLoginResponse)
		case proto.Protocol_MomoTalk_OutLine:
			MomoTalkOutLine(s, syncReq, syncRsp)
			rsp.MomotalkOutlineResponse = syncRsp.(*proto.MomoTalkOutLineResponse)
		case proto.Protocol_Scenario_List:
			ScenarioList(s, syncReq, syncRsp)
			rsp.ScenarioListResponse = syncRsp.(*proto.ScenarioListResponse)
		case proto.Protocol_Shop_GachaRecruitList:
			ShopGachaRecruitList(s, syncReq, syncRsp)
			rsp.ShopGachaRecruitListResponse = syncRsp.(*proto.ShopGachaRecruitListResponse)
		case proto.Protocol_TimeAttackDungeon_Login:
			TimeAttackDungeonLogin(s, syncReq, syncRsp)
			rsp.TimeAttackDungeonLoginResponse = syncRsp.(*proto.TimeAttackDungeonLoginResponse)
		case proto.Protocol_Billing_PurchaseListByYostar:
			BillingPurchaseListByYostar(s, syncReq, syncRsp)
			rsp.BillingPurchaseListByYostarResponse = syncRsp.(*proto.BillingPurchaseListByYostarResponse)
		case proto.Protocol_EventContent_PermanentList:
			EventContentPermanentList(s, syncReq, syncRsp)
			rsp.EventContentPermanentListResponse = syncRsp.(*proto.EventContentPermanentListResponse)
		case proto.Protocol_Attachment_Get:
			AttachmentGet(s, syncReq, syncRsp)
			rsp.AttachmentGetResponse = syncRsp.(*proto.AttachmentGetResponse)
		case proto.Protocol_Attachment_EmblemList:
			AttachmentEmblemList(s, syncReq, syncRsp)
			rsp.AttachmentEmblemListResponse = syncRsp.(*proto.AttachmentEmblemListResponse)
		case proto.Protocol_Sticker_Login:
			StickerLogin(s, syncReq, syncRsp)
			rsp.StickerListResponse = syncRsp.(*proto.StickerLoginResponse)
		case proto.Protocol_MultiFloorRaid_Sync:
			MultiFloorRaidSync(s, syncReq, syncRsp)
			rsp.MultiFloorRaidSyncResponse = syncRsp.(*proto.MultiFloorRaidSyncResponse)
		case proto.Protocol_ContentSweep_MultiSweepPresetList:
			ContentSweepMultiSweepPresetList(s, syncReq, syncRsp)
			rsp.ContentSweepMultiSweepPresetListResponse = syncRsp.(*proto.ContentSweepMultiSweepPresetListResponse)
		default:
			logger.Error("AccountLoginSync 没有处理的数据:%s", cmdId.String())
		}
	}
}

func ContentSaveGet(s *enter.Session, request, response mx.Message) {

}

func ProofTokenSubmit(s *enter.Session, request, response mx.Message) {

}

func AccountSetRepresentCharacterAndComment(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountSetRepresentCharacterAndCommentRequest)
	rsp := response.(*proto.AccountSetRepresentCharacterAndCommentResponse)

	game.SetComment(s, req.Comment)
	game.SetLobbyStudent(s, req.RepresentCharacterServerId)
	rsp.AccountDB = game.GetAccountDB(s)
	rsp.RepresentCharacterDB = game.GetCharacterDB(s, game.GetRepresentCharacterUniqueId(s))
}

func ScenarioAccountStudentChange(s *enter.Session, request, response mx.Message) {

}

func ScenarioLobbyStudentChange(s *enter.Session, request, response mx.Message) {
	// 逆天玩意,纯本地你发服务端干什么?
}

func ToastList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ToastListResponse)

	for _, str := range game.GetToast(s) {
		rsp.ToastDBs = append(rsp.ToastDBs, &proto.ToastDB{
			UniqueId:  0,
			Text:      str,
			ToastId:   str,
			BeginDate: mx.Now(),
			EndDate:   mx.Now().Add(30 * time.Second),
			LifeTime:  3000, // ms
			Delay:     0})
	}
	game.DelToast(s)
}

func ContentSweepMultiSweepPresetList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ContentSweepMultiSweepPresetListResponse)

	rsp.MultiSweepPresetDBs = make([]*proto.MultiSweepPresetDB, 0)
}
