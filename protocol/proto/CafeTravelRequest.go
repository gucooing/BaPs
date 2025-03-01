package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type CafeTravelRequest struct {
	message ProtoMessage
	ResponsePacket

	TargetAccountId          int64
	CurrentVisitingAccountId int64
}

func (x *CafeTravelRequest) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *CafeTravelRequest) ProtoReflect() Message {
	return x
}

func (x *CafeTravelRequest) GetProtocol() int32 {
	return mx.Protocol_Cafe_Travel
}

func (x *CafeTravelRequest) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
