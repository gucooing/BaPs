// Code generated by gucooing DO NOT EDIT.

package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type CafeReceiveCurrencyRequest struct{
	message ProtoMessage
	RequestPacket

    AccountServerId int64
    CafeDBId int64
}

func (x *CafeReceiveCurrencyRequest) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *CafeReceiveCurrencyRequest) ProtoReflect() Message {
	return x
}

func (x *CafeReceiveCurrencyRequest) GetProtocol() int32 {
	return mx.Protocol_Cafe_ReceiveCurrency
}

func (x *CafeReceiveCurrencyRequest) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}

