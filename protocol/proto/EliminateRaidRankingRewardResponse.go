// Code generated by gucooing DO NOT EDIT.

package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type EliminateRaidRankingRewardResponse struct{
	message ProtoMessage
	ResponsePacket

    ReceivedRankingRewardId int64
    ParcelResultDB *ParcelResultDB
}

func (x *EliminateRaidRankingRewardResponse) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *EliminateRaidRankingRewardResponse) ProtoReflect() Message {
	return x
}

func (x *EliminateRaidRankingRewardResponse) GetProtocol() int32 {
	return mx.Protocol_EliminateRaid_RankingReward
}

func (x *EliminateRaidRankingRewardResponse) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}

