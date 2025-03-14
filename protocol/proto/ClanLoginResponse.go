// Code generated by gucooing DO NOT EDIT.

package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type ClanLoginResponse struct{
	message ProtoMessage
	ResponsePacket

    AccountClanDB *ClanDB
    AccountClanMemberDB *ClanMemberDB
    AssistCharacterDBs []*AssistCharacterDB
    ClanAssistSlotDBs []*ClanAssistSlotDB
}

func (x *ClanLoginResponse) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *ClanLoginResponse) ProtoReflect() Message {
	return x
}

func (x *ClanLoginResponse) GetProtocol() int32 {
	return mx.Protocol_Clan_Login
}

func (x *ClanLoginResponse) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}

