package proto

import (
	"github.com/gucooing/BaPs/mx/protocol"
)

type AcademyGetInfoRequest struct {
	protocol.ProtoMessage
	RequestPacket
}

func (x *AcademyGetInfoRequest) ProtoReflect() protocol.Message {
	return x
}

func (x *AcademyGetInfoRequest) GetProtocolKey() uint16 {
	return 24000
}

func (x *AcademyGetInfoRequest) GetProtocolValue() string {
	return "Academy_GetInfo"
}
