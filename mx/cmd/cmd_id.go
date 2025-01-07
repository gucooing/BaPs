package cmd

import (
	"github.com/gucooing/BaPs/mx/proto"
)

func (c *CmdProtoMap) registerAllMessage() {
	c.regMsg(proto.Protocol_Queuing_GetTicket, func() any { return new(proto.QueuingGetTicketRequest) }, true)
	c.regMsg(proto.Protocol_Queuing_GetTicket, func() any { return new(proto.QueuingGetTicketResponse) }, false)
	c.regMsg(proto.Protocol_Account_CheckYostar, func() any { return new(proto.AccountCheckYostarRequest) }, true)
	c.regMsg(proto.Protocol_Account_CheckYostar, func() any { return new(proto.AccountCheckYostarResponse) }, false)
	c.regMsg(proto.Protocol_Account_Auth, func() any { return new(proto.AccountAuthRequest) }, true)
	c.regMsg(proto.Protocol_Account_Auth, func() any { return new(proto.AccountAuthResponse) }, false)
	c.regMsg(proto.Protocol_Account_Nickname, func() any { return new(proto.AccountNicknameRequest) }, true)
	c.regMsg(proto.Protocol_Account_Nickname, func() any { return new(proto.AccountNicknameResponse) }, false)
	c.regMsg(proto.Protocol_ProofToken_RequestQuestion, func() any { return new(proto.ProofTokenRequestQuestionRequest) }, true)
	c.regMsg(proto.Protocol_ProofToken_RequestQuestion, func() any { return new(proto.ProofTokenRequestQuestionResponse) }, false)
	c.regMsg(proto.Protocol_NetworkTime_Sync, func() any { return new(proto.NetworkTimeSyncRequest) }, true)
	c.regMsg(proto.Protocol_NetworkTime_Sync, func() any { return new(proto.NetworkTimeSyncResponse) }, false)
	c.regMsg(proto.Protocol_Academy_GetInfo, func() any { return new(proto.AcademyGetInfoRequest) }, true)
	c.regMsg(proto.Protocol_Academy_GetInfo, func() any { return new(proto.AcademyGetInfoResponse) }, false)
	c.regMsg(proto.Protocol_Account_LoginSync, func() any { return new(proto.AccountLoginSyncRequest) }, true)
	c.regMsg(proto.Protocol_Account_LoginSync, func() any { return new(proto.AccountLoginSyncResponse) }, false)
}
