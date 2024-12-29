package game

import (
	sro "github.com/gucooing/BaPs/common/server_only"
)

func NewYostarGame(accountId int64) *sro.PlayerBin {
	bin := &sro.PlayerBin{
		AccountId: accountId,
	}
	return bin
}
