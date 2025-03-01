package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type AccountCheckYostarRequest struct {
	ProtoMessage
	RequestPacket

	UID              int64
	YostarToken      string
	EnterTicket      string
	PassCookieResult bool
	Cookie           string
}

func (x *AccountCheckYostarRequest) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *AccountCheckYostarRequest) ProtoReflect() Message {
	return x
}

func (x *AccountCheckYostarRequest) GetProtocol() int32 {
	return mx.Protocol_Account_CheckYostar
}

func (x *AccountCheckYostarRequest) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
