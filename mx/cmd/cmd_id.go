package cmd

import (
	"github.com/gucooing/BaPs/mx/proto"
)

func (c *CmdProtoMap) registerAllMessage() {
	c.regMsg(proto.Protocol_Queuing_GetTicket, func() any { return new(proto.QueuingGetTicketRequest) }, true)
	c.regMsg(proto.Protocol_Queuing_GetTicket, func() any { return new(proto.QueuingGetTicketResponse) }, false)
}
