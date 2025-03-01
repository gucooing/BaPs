// Code generated by gucooing DO NOT EDIT.

package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type EquipmentItemTierUpRequest struct{
	message ProtoMessage
	RequestPacket

    TargetEquipmentServerId int64
    ReplaceInfos []*SelectTicketReplaceInfo
}

func (x *EquipmentItemTierUpRequest) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *EquipmentItemTierUpRequest) ProtoReflect() Message {
	return x
}

func (x *EquipmentItemTierUpRequest) GetProtocol() int32 {
	return mx.Protocol_Equipment_TierUp
}

func (x *EquipmentItemTierUpRequest) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}

