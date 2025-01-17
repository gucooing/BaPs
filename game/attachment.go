package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewAttachmentBin() *sro.AttachmentBin {
	info := &sro.AttachmentBin{
		EmblemList: make(map[int64]*sro.EmblemInfo),
	}
	for _, v := range gdconf.GetEmblemExcelCategoryList("Default") {
		info.EmblemList[v.Id] = &sro.EmblemInfo{
			EmblemId:    v.Id,
			ReceiveDate: time.Now().Unix(),
		}
	}

	return info
}

func GetAttachmentBin(s *enter.Session) *sro.AttachmentBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.AttachmentBin == nil {
		bin.AttachmentBin = NewAttachmentBin()
	}
	return bin.AttachmentBin
}

func GetEmblemInfoList(s *enter.Session) map[int64]*sro.EmblemInfo {
	bin := GetAttachmentBin(s)
	if bin == nil {
		return nil
	}
	if bin.EmblemList == nil {
		bin.EmblemList = make(map[int64]*sro.EmblemInfo)
	}
	return bin.EmblemList
}

func UpEmblemInfoList(s *enter.Session, uniqueIds []int64) {
	bin := GetAttachmentBin(s)
	if bin == nil {
		return
	}
	if bin.EmblemList == nil {
		bin.EmblemList = make(map[int64]*sro.EmblemInfo)
	}
	for _, id := range uniqueIds {
		if conf := gdconf.GetEmblemExcel(id); conf != nil {
			bin.EmblemList[id] = &sro.EmblemInfo{
				EmblemId:    id,
				ReceiveDate: time.Now().Unix(),
			}
		} else {

		}
	}
}

func GetAccountAttachmentDB(s *enter.Session) *proto.AccountAttachmentDB {
	if s == nil {
		return nil
	}
	return &proto.AccountAttachmentDB{
		AccountId:      s.AccountServerId,
		EmblemUniqueId: GetEmblemUniqueId(s),
	}
}
