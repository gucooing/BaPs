package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetTimeAttackBin(s *enter.Session) *sro.TimeAttackBin {
	bin := GetDungeonBin(s)
	conf := gdconf.GetCurTimeAttackDungeonSeasonManageExcelTable()
	if bin == nil || conf == nil {
		return nil
	}
	if bin.TimeAttackBin == nil {
		bin.TimeAttackBin = &sro.TimeAttackBin{}
	}
	if conf.Id != bin.TimeAttackBin.SeasonId {
		bin.TimeAttackBin = &sro.TimeAttackBin{
			SeasonId: conf.Id,
		}
	}
	info := bin.TimeAttackBin
	if alg.GetDay4().After(time.Unix(info.LastUpTime, 0)) {
		info.PreviousRoom = 0
		info.TimeAttackRoomList = make(map[int64]*sro.TimeAttackRoom)
		info.LastUpTime = time.Now().Unix()
	}

	return info
}

func GetTimeAttackRoomList(s *enter.Session) map[int64]*sro.TimeAttackRoom {
	bin := GetTimeAttackBin(s)
	if bin == nil {
		return nil
	}
	if bin.TimeAttackRoomList == nil {
		bin.TimeAttackRoomList = make(map[int64]*sro.TimeAttackRoom)
	}
	return bin.TimeAttackRoomList
}

func GetTimeAttackRoom(s *enter.Session, roomId int64) *sro.TimeAttackRoom {
	list := GetTimeAttackRoomList(s)
	if list == nil {
		return nil
	}
	return list[roomId]
}

func AddCurTimeAttackRoom(s *enter.Session, isPractice bool) {
	bin := GetTimeAttackBin(s)
	if bin == nil {
		return
	}
	if bin.TimeAttackRoomList == nil {
		bin.TimeAttackRoomList = make(map[int64]*sro.TimeAttackRoom)
	}
	roomId := GetServerId(s)
	info := &sro.TimeAttackRoom{
		SeasonId:   bin.SeasonId,
		RoomId:     roomId,
		StartTime:  time.Now().Unix(),
		IsPractice: isPractice,
		RewardTime: 0,
	}
	bin.CurRoom = roomId
	bin.TimeAttackRoomList[roomId] = info
}

func GetTimeAttackDungeonRoomDB(s *enter.Session, roomId int64) *proto.TimeAttackDungeonRoomDB {
	bin := GetTimeAttackRoom(s, roomId)
	if bin == nil {
		return nil
	}
	info := &proto.TimeAttackDungeonRoomDB{
		AccountId:         s.AccountServerId,
		SeasonId:          bin.SeasonId,                                       // 赛季id
		RoomId:            bin.RoomId,                                         // 房间id
		CreateDate:        mx.Unix(bin.StartTime, 0),                          // 战斗开始时间
		RewardDate:        mx.Unix(bin.RewardTime, 0),                         // 奖励领取时间
		IsPractice:        bin.IsPractice,                                     // 是否实践
		SweepHistoryDates: make([]time.Time, 0),                               // 扫荡记录
		BattleHistoryDBs:  make([]*proto.TimeAttackDungeonBattleHistoryDB, 0), // 成功战斗记录
	}
	return info
}

func GetTimeAttackDungeonRoomDBs(s *enter.Session) map[int64]*proto.TimeAttackDungeonRoomDB {
	list := make(map[int64]*proto.TimeAttackDungeonRoomDB, 0)
	for id, _ := range GetTimeAttackRoomList(s) {
		info := GetTimeAttackDungeonRoomDB(s, id)
		if info == nil {
			continue
		}
		list[id] = info
	}
	return list
}

func GetCurTimeAttackDungeonRoomDB(s *enter.Session) *proto.TimeAttackDungeonRoomDB {
	bin := GetTimeAttackBin(s)
	if bin == nil {
		return nil
	}
	return GetTimeAttackDungeonRoomDB(s, bin.CurRoom)
}
