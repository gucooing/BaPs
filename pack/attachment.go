package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func AttachmentGet(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.AttachmentGetResponse)

	rsp.AccountAttachmentDB = game.GetAccountAttachmentDB(s)
}

func AttachmentEmblemList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.AttachmentEmblemListResponse)

	rsp.EmblemDBs = make([]*proto.EmblemDB, 0)
	for _, v := range game.GetEmblemInfoList(s) {
		rsp.EmblemDBs = append(rsp.EmblemDBs, &proto.EmblemDB{
			Type:        proto.ParcelType_IdCardBackground,
			UniqueId:    v.EmblemId,
			ReceiveDate: mx.Unix(v.ReceiveDate, 0),
			ParcelInfos: make([]*proto.ParcelInfo, 0),
		})
	}
}

func AttachmentEmblemAcquire(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AttachmentEmblemAcquireRequest)
	rsp := response.(*proto.AttachmentEmblemAcquireResponse)

	game.UpEmblemInfoList(s, req.UniqueIds)

	rsp.EmblemDBs = make([]*proto.EmblemDB, 0)
	for _, v := range game.GetEmblemInfoList(s) {
		rsp.EmblemDBs = append(rsp.EmblemDBs, &proto.EmblemDB{
			Type:        proto.ParcelType_IdCardBackground,
			UniqueId:    v.EmblemId,
			ReceiveDate: mx.Unix(v.ReceiveDate, 0),
			ParcelInfos: make([]*proto.ParcelInfo, 0),
		})
	}
}

func AttachmentEmblemAttach(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AttachmentEmblemAttachRequest)
	rsp := response.(*proto.AttachmentEmblemAttachResponse)

	game.SetEmblemUniqueId(s, req.UniqueId)
	rsp.AttachmentDB = game.GetAccountAttachmentDB(s)
}
