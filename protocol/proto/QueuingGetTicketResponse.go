package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type QueuingGetTicketResponse struct {
	message ProtoMessage
	ResponsePacket

	WaitingTicket          string
	EnterTicket            string
	TicketSequence         int64
	AllowedSequence        int64
	RequiredSecondsPerUser float64
	Birth                  string
	ServerSeed             string
}

func (x *QueuingGetTicketResponse) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *QueuingGetTicketResponse) ProtoReflect() Message {
	return x
}

func (x *QueuingGetTicketResponse) GetProtocol() int32 {
	return mx.Protocol_Queuing_GetTicket
}

func (x *QueuingGetTicketResponse) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
