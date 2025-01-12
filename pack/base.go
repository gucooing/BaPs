package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func MemoryLobbyList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MemoryLobbyListResponse)

	rsp.MemoryLobbyDBs = make([]*proto.MemoryLobbyDB, 0)
}

func CampaignList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CampaignListResponse)

	rsp.StageHistoryDBs = make([]*proto.CampaignStageHistoryDB, 0)
}

func TimeAttackDungeonLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.TimeAttackDungeonLoginResponse)

	rsp.PreviousRoomDB = &proto.TimeAttackDungeonRoomDB{
		AccountId:         0,
		SeasonId:          0,
		RoomId:            0,
		CreateDate:        time.Time{},
		RewardDate:        time.Time{},
		IsPractice:        false,
		SweepHistoryDates: nil,
		BattleHistoryDBs:  nil,
	}
}

func BillingPurchaseListByYostar(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.BillingPurchaseListByYostarResponse)

	rsp.CountList = make([]*proto.PurchaseCountDB, 0)
	rsp.OrderList = make([]*proto.PurchaseOrderDB, 0)
	rsp.MonthlyProductList = make([]*proto.MonthlyProductPurchaseDB, 0)
	rsp.BlockedProductDBs = make([]*proto.BlockedProductDB, 0)
}

func EventContentPermanentList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.EventContentPermanentListResponse)

	rsp.PermanentDBs = make([]*proto.EventContentPermanentDB, 0)
	for _, id := range []int64{900801, 900802, 900803, 900804, 900805, 900806, 900808, 900809,
		900810, 900812, 900813, 900814, 900815, 900816, 900817, 900818, 900825, 900701} {
		rsp.PermanentDBs = append(rsp.PermanentDBs, &proto.EventContentPermanentDB{
			EventContentId:            id,
			IsStageAllClear:           false,
			IsReceivedCharacterReward: false,
		})
	}
}

func AttachmentGet(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.AttachmentGetResponse)

	rsp.AccountAttachmentDB = &proto.AccountAttachmentDB{
		AccountId:      s.AccountServerId,
		EmblemUniqueId: 0,
	}
}

func AttachmentEmblemList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.AttachmentEmblemListResponse)

	rsp.EmblemDBs = make([]*proto.EmblemDB, 0)
	for _, id := range []int64{1, 2, 3, 4, 5} {
		rsp.EmblemDBs = append(rsp.EmblemDBs, &proto.EmblemDB{
			Type:        proto.ParcelType_IdCardBackground,
			UniqueId:    id,
			ReceiveDate: time.Now(),
			ParcelInfos: make([]*proto.ParcelInfo, 0),
		})
	}
}

func StickerLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.StickerLoginResponse)

	rsp.StickerBookDB = &proto.StickerBookDB{
		AccountId:        s.AccountServerId,
		UnusedStickerDBs: make([]*proto.StickerDB, 0),
		UsedStickerDBs:   make([]*proto.StickerDB, 0),
	}
}
