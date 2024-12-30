package game

import (
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/alg"
)

func NewYostarGame(accountId int64) *sro.PlayerBin {
	bin := &sro.PlayerBin{
		BaseBin: &sro.BasePlayer{
			AccountId: accountId,
			Level:     1,
			Nickname:  "",
			CreateDate: alg.GetTimestampProto(time.Now()),
		},
	}
	return bin
}
