// Code generated by gucooing DO NOT EDIT.

package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type ShopBeforehandGachaRunRequest struct{
	message ProtoMessage
	RequestPacket

    ShopUniqueId int64
    GoodsId int64
}

func (x *ShopBeforehandGachaRunRequest) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *ShopBeforehandGachaRunRequest) ProtoReflect() Message {
	return x
}

func (x *ShopBeforehandGachaRunRequest) GetProtocol() int32 {
	return mx.Protocol_Shop_BeforehandGachaRun
}

func (x *ShopBeforehandGachaRunRequest) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}

