package game

import (
	"math"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func GetDBId() int64 {
	return 123456
}

func GetServerId(s *enter.Session) int64 {
	bin := GetPlayerBin(s)
	if bin == nil {
		return 0
	}
	if bin.ServerId == math.MaxInt64 {
		logger.Warn("[UID:%v]玩家唯一计数器达到最大值:%v", s.AccountServerId, bin.ServerId)
	}
	if bin.ServerId < 1e8 {
		bin.ServerId = 1e8
	}
	defer func() {
		bin.ServerId++
	}()
	return bin.ServerId
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
