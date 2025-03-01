package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type AccountCheckYostarResponse struct {
	message ProtoMessage
	ResponsePacket

	ResultState  int
	ResultMessag string
	Birth        string
}

func (x *AccountCheckYostarResponse) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *AccountCheckYostarResponse) ProtoReflect() Message {
	return x
}

func (x *AccountCheckYostarResponse) GetProtocol() int32 {
	return mx.Protocol_Account_CheckYostar
}

func (x *AccountCheckYostarResponse) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
