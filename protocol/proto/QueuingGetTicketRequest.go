package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type QueuingGetTicketRequest struct {
	ProtoMessage
	RequestPacket

	YostarUID       int64
	YostarToken     string
	MakeStandby     bool
	PassCheck       bool
	PassCheckYostar bool
	WaitingTicket   string
	ClientVersion   string
}

func (x *QueuingGetTicketRequest) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *QueuingGetTicketRequest) ProtoReflect() Message {
	return x
}

func (x *QueuingGetTicketRequest) GetProtocol() int32 {
	return mx.Protocol_Queuing_GetTicket
}

func (x *QueuingGetTicketRequest) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
