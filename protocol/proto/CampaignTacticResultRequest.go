// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type CampaignTacticResultRequest struct {
	*RequestPacket

	PassCheckCharacter bool
	Summary            *BattleSummary
	Hand               *SkillCardHand
	SkipSummary        *TacticSkipSummary
}

func (x *CampaignTacticResultRequest) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *CampaignTacticResultRequest) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.RequestPacket = packet.(*RequestPacket)
}
