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
		if ok {
			flagS += flag
		}
	}
	return flagS
}

func GetMissionProgressDBs(s *enter.Session) []*proto.MissionProgressDB {
	list := make([]*proto.MissionProgressDB, 0)
	for _, bin := range s.GetMissionSync() {
		info := GetMissionProgressDB(s, bin.MissionId)
		if info != nil {
			list = append(list, info)
		}
	}
	s.NewMissionSync()
	return list
}

func GetToast(s *enter.Session) []string {
	if s == nil {
		return nil
	}
	if s.Toast == nil {
		s.Toast = make([]string, 0)
	}
	return s.Toast
}

func AddToast(s *enter.Session, toast string) {
	if s == nil {
		return
	}
	if s.Toast == nil {
		s.Toast = make([]string, 0)
	}
	s.Toast = append(s.Toast, toast)
	SetServerNotification(s, proto.ServerNotificationFlag_NewToastDetected, true)
}

func DelToast(s *enter.Session) {
	if s == nil {
		return
	}
	s.Toast = make([]string, 0)
	SetServerNotification(s, proto.ServerNotificationFlag_NewToastDetected, false)
}
