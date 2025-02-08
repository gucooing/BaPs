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

func GetServerNotification(s *enter.Session) int32 {
	flagS := int32(0)
	for flag, ok := range GetActions(s) {
		if ok {
			flagS += int32(flag)
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
