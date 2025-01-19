package game

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetActions(s *enter.Session) map[proto.ServerNotificationFlag]bool {
	if s == nil {
		return nil
	}
	if s.Actions == nil {
		s.Actions = make(map[proto.ServerNotificationFlag]bool)
	}
	return s.Actions
}

func SetServerNotification(s *enter.Session, flag proto.ServerNotificationFlag, ok bool) {
	if s == nil {
		return
	}
	if s.Actions == nil {
		s.Actions = make(map[proto.ServerNotificationFlag]bool)
	}
	s.Actions[flag] = ok
}

func GetServerNotification(s *enter.Session) proto.ServerNotificationFlag {
	flagS := proto.ServerNotificationFlag_None
	for flag, ok := range GetActions(s) {
		if ok && flagS < flag {
			flagS = flag
		}
	}
	return flagS
}
