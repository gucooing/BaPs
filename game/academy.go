package game

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetAcademyDB(s *enter.Session) *proto.AcademyDB {
	return &proto.AcademyDB{}
}

func GetAcademyLocationDBs(s *enter.Session) []*proto.AcademyLocationDB {
	return make([]*proto.AcademyLocationDB, 0)
}
