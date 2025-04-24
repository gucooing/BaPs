package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func MailCheck(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.MailCheckResponse)

	rsp.Count = game.GetMailCheckCount(s) // 未领取数量
	if rsp.Count > 0 {
		game.SetServerNotification(s, proto.ServerNotificationFlag_HasUnreadMail, true)
	} else {
		game.SetServerNotification(s, proto.ServerNotificationFlag_HasUnreadMail, false)
	}
}

func MailList(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.MailListRequest)
	rsp := response.(*proto.MailListResponse)

	rsp.MailDBs = game.GetMailDBs(s, req.IsReadMail)
	game.SetServerNotification(s, proto.ServerNotificationFlag_HasUnreadMail, false)
}

func MailReceive(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.MailReceiveRequest)
	rsp := response.(*proto.MailReceiveResponse)

	rsp.MailServerIds = make([]int64, 0)
	result := make([]*game.ParcelResult, 0)
	readParcelNum := 0
	for _, mailServerId := range req.MailServerIds {
		bin := game.GetMailInfo(s, mailServerId)
		if bin == nil || bin.IsRead {
			logger.Debug("[UID:%v]没有找到或者已领取邮件Id:%v", s.AccountServerId, mailServerId)
			continue
		}
		result = append(result, game.GetMailParcelResultList(bin.ParcelInfoList)...)
		bin.ReceiptDate = time.Now().Unix()
		bin.IsRead = true
		rsp.MailServerIds = append(rsp.MailServerIds, mailServerId)
		readParcelNum += len(bin.ParcelInfoList)
		if readParcelNum >= game.MaxMailParcelNum {
			goto ty
		}
	}
ty:
	rsp.ParcelResultDB = game.ParcelResultDB(s, result)
}
