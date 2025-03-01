package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type ShopBuyGacha2Request struct {
	message ProtoMessage
	ResponsePacket

	ShopBuyGachaRequest
}

func (x *ShopBuyGacha2Request) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *ShopBuyGacha2Request) ProtoReflect() Message {
	return x
}

func (x *ShopBuyGacha2Request) GetProtocol() int32 {
	return mx.Protocol_Shop_BuyGacha2
}

func (x *ShopBuyGacha2Request) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
