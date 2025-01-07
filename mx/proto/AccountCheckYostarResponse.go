package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/mx"
)

type AccountCheckYostarResponse struct {
	message mx.ProtoMessage
	ResponsePacket

	ResultState  int
	ResultMessag string
	Birth        string
}

func (x *AccountCheckYostarResponse) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *AccountCheckYostarResponse) ProtoReflect() mx.Message {
	return x
}

func (x *AccountCheckYostarResponse) GetProtocol() int32 {
	return Protocol_Account_CheckYostar
}

func (x *AccountCheckYostarResponse) SetSessionKey(base *mx.BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
