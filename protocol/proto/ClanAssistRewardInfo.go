// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type ClanAssistRewardInfo struct {
	CharacterDBId            int64
	DeployDate               mx.MxTime
	RentCount                int64
	CumultativeRewardParcels []*ParcelInfo
	RentRewardParcels        []*ParcelInfo
}

func (x *ClanAssistRewardInfo) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
