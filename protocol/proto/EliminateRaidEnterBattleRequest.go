// Code generated by gucooing DO NOT EDIT.

package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type EliminateRaidEnterBattleRequest struct{
	message ProtoMessage
	RequestPacket

    RaidServerId int64
    RaidUniqueId int64
    IsPractice bool
    EchelonId int64
    AssistUseInfo *ClanAssistUseInfo
}

func (x *EliminateRaidEnterBattleRequest) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *EliminateRaidEnterBattleRequest) ProtoReflect() Message {
	return x
}

func (x *EliminateRaidEnterBattleRequest) GetProtocol() int32 {
	return mx.Protocol_EliminateRaid_EnterBattle
}

func (x *EliminateRaidEnterBattleRequest) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}

