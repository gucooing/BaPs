// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
)

type ITBGItemInfo struct {
	UniqueId           int64
	ItemType           TBGItemType
	ItemEffectType     TBGItemEffectType
	ItemParameter      int32
	EncounterCount     int32
	LocalizeEtcId      string
	Icon               string
	BuffIcon           string
	DiceEffectAniClip  string
	BuffIconHUDVisible bool
}

func (x *ITBGItemInfo) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
