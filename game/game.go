package game

import (
	"math"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetDBId() int64 {
	return 123456
}

func GetServerId(s *enter.Session) int64 {
	if s == nil ||
		s.PlayerBin == nil {
		return 0
	}
	if s.PlayerBin.ServerId == math.MaxInt64 {
		logger.Warn("[UID:%v]玩家唯一计数器达到最大值:%v", s.AccountServerId, s.PlayerBin.ServerId)
	}
	s.PlayerBin.ServerId++
	return s.PlayerBin.ServerId
}

func GetServerTime() int64 {
	return (time.Now().Unix() * 10000000) + 621356292000000000
}

func GetServerTimeTick() int64 {
	return time.Now().UnixNano()/100 + 621356292000000000
}

func GetServerNotification(s *enter.Session) proto.ServerNotificationFlag {
	return proto.ServerNotificationFlag_None
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
