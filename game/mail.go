package game

import (
	dbstruct "github.com/gucooing/BaPs/db/struct"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

const (
	MaxMailParcelNum = 700
)

func NewMailBin(s *enter.Session) *sro.MailBin {
	return &sro.MailBin{}
}

func GetMailBin(s *enter.Session) *sro.MailBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.MailBin == nil {
		bin.MailBin = NewMailBin(s)
	}
	return bin.MailBin
}

// 检查方法
func GetMailCheckCount(s *enter.Session) int64 {
	bin := GetMailBin(s)
	if bin == nil {
		return 0
	}
	if bin.YostarMail == nil {
		bin.YostarMail = make(map[int64]bool)
	}
	if bin.MailInfoList == nil {
		bin.MailInfoList = make(map[int64]*sro.MailInfo)
	}
	var count int64 = 0
	// 检查是否有全局邮件没有添加到玩家数据
	for v, k := range enter.GetYostarMail() {
		if ok := bin.YostarMail[v]; !ok {
			id := GetServerId(s)
			bin.MailInfoList[id] = &sro.MailInfo{
				ServerId:       id,
				Sender:         k.Sender,
				Comment:        k.Comment,
				SendDate:       k.SendDate.Time.Unix(),
				ExpireDate:     k.ExpireDate.Time.Unix(),
				ParcelInfoList: MailParcelInfoJsonToProtobuf(k.ParcelInfoList),
				IsRead:         false,
			}
			bin.YostarMail[v] = true
		}
	}
	// 检查未读邮件
	for serverId, info := range bin.MailInfoList {
		// 删除超时邮件
		if time.Now().After(time.Unix(info.ExpireDate, 0)) {
			delete(bin.MailInfoList, serverId)
			continue
		}
		if !info.IsRead {
			count++
		}
	}
	return count
}

func MailParcelInfoJsonToProtobuf(bin []*dbstruct.ParcelInfo) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, len(bin))
	for i, v := range bin {
		list[i] = &sro.ParcelInfo{
			Type: v.Type,
			Id:   v.Id,
			Num:  v.Num,
		}
	}
	return list
}

func GetMailInfoList(s *enter.Session) map[int64]*sro.MailInfo {
	bin := GetMailBin(s)
	if bin == nil {
		return nil
	}
	if bin.MailInfoList == nil {
		bin.MailInfoList = make(map[int64]*sro.MailInfo)
	}
	return bin.MailInfoList
}

func GetMailInfo(s *enter.Session, id int64) *sro.MailInfo {
	bin := GetMailInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[id]
}

func AddMail(s *enter.Session, info *sro.MailInfo) bool {
	bin := GetMailBin(s)
	if bin == nil {
		return false
	}
	if bin.MailInfoList == nil {
		bin.MailInfoList = make(map[int64]*sro.MailInfo)
	}
	serverId := GetServerId(s)
	info.ServerId = serverId
	bin.MailInfoList[serverId] = info
	SetServerNotification(s, proto.ServerNotificationFlag_HasUnreadMail, true)
	return true
}

var SystemMail = map[string]string{
	"Attendance":     "请查收每日登录奖励",
	"MaxActionPoint": "超出的体力",
}

func AddMailBySystem(s *enter.Session, mailType string, parcelInfoList []*sro.ParcelInfo) bool {
	bin := GetMailBin(s)
	if bin == nil {
		return false
	}
	if bin.MailInfoList == nil {
		bin.MailInfoList = make(map[int64]*sro.MailInfo)
	}
	serverId := GetServerId(s)
	mail := &sro.MailInfo{
		Sender:         "gucooing",
		ServerId:       serverId,
		Comment:        SystemMail[mailType],
		SendDate:       time.Now().Unix(),
		ExpireDate:     time.Now().Add(7 * 24 * time.Hour).Unix(),
		ParcelInfoList: parcelInfoList,
	}

	bin.MailInfoList[serverId] = mail
	SetServerNotification(s, proto.ServerNotificationFlag_HasUnreadMail, true)
	return true
}

func DelMail(s *enter.Session, id int64) bool {
	bin := GetMailBin(s)
	if bin == nil {
		return false
	}
	if _, ok := bin.MailInfoList[id]; ok {
		delete(bin.MailInfoList, id)
		return true
	}
	return false
}

func GetMailDBs(s *enter.Session, IsReadMail bool) []*proto.MailDB {
	list := make([]*proto.MailDB, 0)
	for _, bin := range GetMailInfoList(s) {
		if bin.IsRead == IsReadMail {
			info := &proto.MailDB{
				ServerId:          bin.ServerId,
				AccountServerId:   s.AccountServerId,
				Type:              0,
				UniqueId:          1,
				Sender:            bin.Sender,
				Comment:           bin.Comment,
				SendDate:          mx.Unix(bin.SendDate, 0),
				ReceiptDate:       mx.Unix(bin.ReceiptDate, 0),
				ExpireDate:        mx.Unix(bin.ExpireDate, 0),
				ParcelInfos:       make([]*proto.ParcelInfo, 0),
				RemainParcelInfos: make([]*proto.ParcelInfo, 0),
			}
			for _, parcelInfo := range bin.ParcelInfoList {
				info.ParcelInfos = append(info.ParcelInfos, &proto.ParcelInfo{
					Key: &proto.ParcelKeyPair{
						Type: proto.ParcelType(parcelInfo.Type),
						Id:   parcelInfo.Id,
					},
					Amount: parcelInfo.Num,
					Multiplier: &proto.BasisPoint{
						RawValue: 10000,
					},
					Probability: &proto.BasisPoint{
						RawValue: 10000,
					},
				})
			}
			list = append(list, info)
		}
	}

	return list
}

func GetMailParcelResultList(bin []*sro.ParcelInfo) []*ParcelResult {
	list := make([]*ParcelResult, 0)
	for _, info := range bin {
		list = append(list, &ParcelResult{
			ParcelType: proto.ParcelType(info.Type),
			ParcelId:   info.Id,
			Amount:     info.Num,
		})
	}
	return list
}
