package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type CafeTravelResponse struct {
	message ProtoMessage
	ResponsePacket

	FriendDB *FriendDB
	CafeDBs  []*CafeDB
}

func (x *CafeTravelResponse) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *CafeTravelResponse) ProtoReflect() Message {
	return x
}

func (x *CafeTravelResponse) GetProtocol() int32 {
	return mx.Protocol_Cafe_Travel
}

func (x *CafeTravelResponse) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
