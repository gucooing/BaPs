package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
)

func GetServerId(s *enter.Session) int64 {
	sw := alg.GetSnow()
	return sw.GenId()
}

func GetServerTime() int64 {
	return (time.Now().Add(-1*time.Hour).Unix() * 10000000) + 621356292000000000
}

func GetServerTimeTick() int64 {
	return time.Now().Add(-1*time.Hour).UnixNano()/100 + 621356292000000000
}

func NewYostarGame(accountId int64) *sro.PlayerBin {
	bin := &sro.PlayerBin{
		BaseBin: &sro.BasePlayer{
			AccountId:  accountId,
			Level:      1,
			Nickname:   "",
			CreateDate: time.Now().Unix(),
		},
	}
	return bin
}

func GetPlayerBin(s *enter.Session) *sro.PlayerBin {
	if s == nil ||
		s.PlayerBin == nil {
		logger.Error("数据损坏")
		return nil
	}
	return s.PlayerBin
}
