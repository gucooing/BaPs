package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type ShopBuyGacha3Request struct {
	message ProtoMessage
	ResponsePacket

	ShopBuyGacha2Request
	FreeRecruitId int64
	Cost          *ParcelCost
}

func (x *ShopBuyGacha3Request) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *ShopBuyGacha3Request) ProtoReflect() Message {
	return x
}

func (x *ShopBuyGacha3Request) GetProtocol() int32 {
	return mx.Protocol_Shop_BuyGacha3
}

func (x *ShopBuyGacha3Request) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
