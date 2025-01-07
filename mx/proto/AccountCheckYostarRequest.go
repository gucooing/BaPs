package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/mx"
)

type AccountCheckYostarRequest struct {
	mx.ProtoMessage
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

func (x *AccountCheckYostarRequest) ProtoReflect() mx.Message {
	return x
}

func (x *AccountCheckYostarRequest) GetProtocol() int32 {
	return Protocol_Account_CheckYostar
}

func (x *AccountCheckYostarRequest) SetSessionKey(base *mx.BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
