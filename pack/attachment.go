package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/protocol/proto"
)

func AttachmentGet(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.AttachmentGetResponse)

	rsp.AccountAttachmentDB = game.GetAccountAttachmentDB(s)
}

func AttachmentEmblemList(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.AttachmentEmblemListResponse)

	rsp.EmblemDBs = game.GetEmblemDBs(s)
}

func AttachmentEmblemAcquire(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.AttachmentEmblemAcquireRequest)
	rsp := response.(*proto.AttachmentEmblemAcquireResponse)

	game.UpEmblemInfoList(s, req.UniqueIds)
	rsp.EmblemDBs = game.GetEmblemDBs(s)
}

func AttachmentEmblemAttach(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.AttachmentEmblemAttachRequest)
	rsp := response.(*proto.AttachmentEmblemAttachResponse)

	game.SetEmblemUniqueId(s, req.UniqueId)
	rsp.AttachmentDB = game.GetAccountAttachmentDB(s)
}
