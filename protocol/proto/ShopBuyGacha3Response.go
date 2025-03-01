package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type ShopBuyGacha3Response struct {
	message ProtoMessage
	ResponsePacket

	ShopBuyGacha2Response
	FreeRecruitHistoryDB *ShopFreeRecruitHistoryDB
}

func (x *ShopBuyGacha3Response) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *ShopBuyGacha3Response) ProtoReflect() Message {
	return x
}

func (x *ShopBuyGacha3Response) GetProtocol() int32 {
	return mx.Protocol_Shop_BuyGacha3
}

func (x *ShopBuyGacha3Response) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}
