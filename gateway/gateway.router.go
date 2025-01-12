package gateway

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/cmd"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pack"
	"github.com/gucooing/BaPs/pkg/logger"
)

type handlerFunc func(s *enter.Session, request, response mx.Message)

func (g *Gateway) newFuncRouteMap() {
	g.funcRouteMap = map[int32]handlerFunc{
		proto.Protocol_Account_CheckYostar:               g.AccountCheckYostar,                  // 验证EnterTicket
		proto.Protocol_Account_Auth:                      pack.AccountAuth,                      // 账号验证
		proto.Protocol_Account_Nickname:                  pack.AccountNickname,                  // 设置昵称
		proto.Protocol_ProofToken_RequestQuestion:        pack.ProofTokenRequestQuestion,        // 验证登录token
		proto.Protocol_NetworkTime_Sync:                  pack.NetworkTimeSync,                  // 同步时间
		proto.Protocol_Academy_GetInfo:                   pack.AcademyGetInfo,                   // 获取学院信息
		proto.Protocol_Account_LoginSync:                 pack.AccountLoginSync,                 // 同步账号信息
		proto.Protocol_Cafe_Get:                          pack.CafeGetInfo,                      // 获取咖啡馆信息
		proto.Protocol_Account_CurrencySync:              pack.AccountCurrencySync,              // 同步账号货币
		proto.Protocol_Character_List:                    pack.CharacterList,                    // 获取角色列表
		proto.Protocol_Equipment_List:                    pack.EquipmentList,                    // 获取装备信息
		proto.Protocol_CharacterGear_List:                pack.CharacterGearList,                // 获取角色??
		proto.Protocol_Echelon_List:                      pack.EchelonList,                      // ???
		proto.Protocol_MemoryLobby_List:                  pack.MemoryLobbyList,                  // 获取记忆大厅列表
		proto.Protocol_Campaign_List:                     pack.CampaignList,                     // 获取营销??
		proto.Protocol_Arena_Login:                       pack.ArenaLogin,                       // 登录获取竞技场信息
		proto.Protocol_Raid_Login:                        pack.RaidLogin,                        // 登录获取总力站信息
		proto.Protocol_EliminateRaid_Login:               pack.EliminateRaidLogin,               // 登录获取制约解除决战信息
		proto.Protocol_Craft_List:                        pack.CraftInfoList,                    // 获取制造信息
		proto.Protocol_Clan_Login:                        pack.ClanLogin,                        // 登录获取社团信息
		proto.Protocol_MomoTalk_OutLine:                  pack.MomoTalkOutLine,                  // 获取MomoTalk信息
		proto.Protocol_Scenario_List:                     pack.ScenarioList,                     // 获取场景剧情信息（应该是MomoTalk的东西
		proto.Protocol_Shop_GachaRecruitList:             pack.ShopGachaRecruitList,             // 获取商店卡池??
		proto.Protocol_TimeAttackDungeon_Login:           pack.TimeAttackDungeonLogin,           // 登录获取限时战斗信息??
		proto.Protocol_Billing_PurchaseListByYostar:      pack.BillingPurchaseListByYostar,      // Yostar采购清单??
		proto.Protocol_EventContent_PermanentList:        pack.EventContentPermanentList,        // 获取永久时间列表?
		proto.Protocol_Attachment_Get:                    pack.AttachmentGet,                    // 获取附件??
		proto.Protocol_Attachment_EmblemList:             pack.AttachmentEmblemList,             // 获取附件徽章列表??
		proto.Protocol_Sticker_Login:                     pack.StickerLogin,                     // 登录获取贴纸信息??
		proto.Protocol_MultiFloorRaid_Sync:               pack.MultiFloorRaidSync,               // 制约解除决战信息同步??
		proto.Protocol_ContentSweep_MultiSweepPresetList: pack.ContentSweepMultiSweepPresetList, // ????
		proto.Protocol_Item_List:                         pack.ItemList,                         // 获取背包物品
		proto.Protocol_ContentSave_Get:                   pack.ContentSaveGet,                   // ???
		proto.Protocol_Shop_BeforehandGachaGet:           pack.ShopBeforehandGachaGet,           // ???
		proto.Protocol_ProofToken_Submit:                 pack.ProofTokenSubmit,                 // 密钥提交
		proto.Protocol_Account_GetTutorial:               pack.AccountGetTutorial,               // 获取教程
		proto.Protocol_Mission_List:                      pack.MissionList,                      // 获取剧情信息
		proto.Protocol_Mission_GuideMissionSeasonList:    pack.GuideMissionSeasonList,           // 获取指南任务信息
		proto.Protocol_Toast_List:                        pack.ToastList,
	}
}

func (g *Gateway) registerMessage(c *gin.Context, request mx.Message, base *mx.BasePacket) {
	// panic捕获
	defer func() {
		if err := recover(); err != nil {
			errBestHTTP(c)
			logger.Error("@LogTag(player_panic)cmdId:%s json:%s\nerr:%s\nstack:%s", cmd.Get().GetCmdNameByCmdId(request.GetProtocol()),
				request.String(), err, logger.Stack())
			return
		}
	}()
	handler, ok := g.funcRouteMap[request.GetProtocol()]
	if !ok {
		errBestHTTP(c)
		logPlayerMsg(NoRoute, request)
		return
	}
	response := cmd.Get().GetResponsePacketByCmdId(request.GetProtocol())
	if response == nil {
		errBestHTTP(c)
		logger.Debug("response unknown cmd id: %v\n", cmd.Get().GetCmdNameByCmdId(request.GetProtocol()))
		return
	}
	sessionKey := base.GetSessionKey()
	var s *enter.Session
	if sessionKey == nil &&
		request.GetProtocol() != proto.Protocol_Account_CheckYostar {
		errBestHTTP(c)
		logger.Debug("get request sessionKey nil")
		return
	} else if request.GetProtocol() != proto.Protocol_Account_CheckYostar {
		s = enter.GetSessionBySessionKey(sessionKey)
		if s == nil {
			errBestHTTP(c)
			logger.Debug("get session nil")
			return
		}
		s.EndTime = time.Now().Add(30 * time.Minute)
	}
	response.SetSessionKey(base)
	handler(s, request, response)
	logPlayerMsg(Client, request)
	logPlayerMsg(Server, response)
	g.send(c, response)
	return
}

// 62135629200 * 10000000 + curTime

const (
	Client  = 1
	Server  = 2
	NoRoute = 3
)

func logPlayerMsg(logType int, msg mx.Message) {
	if _, ok := config.GetBlackCmd()[proto.Protocol(msg.GetProtocol()).String()]; ok ||
		!config.GetIsLogMsgPlayer() {
		return
	}
	var a string
	switch logType {
	case Client:
		a = "@LogTag(player_msg)@ gateway c--->s cmd id:"
	case Server:
		a = "@LogTag(player_msg)@ gateway s--->c cmd id:"
	case NoRoute:
		a = "@LogTag(player_no_route)@ c --> s no route for msg, cmd id:"
	}
	b, _ := json.MarshalIndent(msg, "", "  ")

	logger.Debug("%s%v :%s", a, msg.GetProtocol(), string(b))
}
