// Code generated by gucooing DO NOT EDIT.

package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type EchelonPresetSaveResponse struct{
	message ProtoMessage
	ResponsePacket

    PresetDB *EchelonPresetDB
}

func (x *EchelonPresetSaveResponse) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *EchelonPresetSaveResponse) ProtoReflect() Message {
	return x
}

func (x *EchelonPresetSaveResponse) GetProtocol() int32 {
	return mx.Protocol_Echelon_PresetSave
}

func (x *EchelonPresetSaveResponse) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}

