package game

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
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
	if alg.GetLastDayH(4).After(time.Unix(info.LastUpTime, 0)) {
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
		StartTime:  time.Now().Add(1 * time.Hour).Unix(),
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
		AccountId:  s.AccountServerId,
		SeasonId:   bin.SeasonId,              // 赛季id
		RoomId:     bin.RoomId,                // 房间id
		CreateDate: mx.Unix(bin.StartTime, 0), // 战斗结束时间
		IsPractice: bin.IsPractice,            // 是否实践
		// RewardDate: mx.Unix(bin.RewardTime, 0), // 奖励领取时间
		SweepHistoryDates: make([]mx.MxTime, 0),                               // 扫荡记录
		BattleHistoryDBs:  make([]*proto.TimeAttackDungeonBattleHistoryDB, 0), // 成功战斗记录
	}
	if bin.RewardTime != 0 {
		info.RewardDate = mx.Unix(bin.RewardTime, 0)
	}
	for _, battleInfo := range bin.BattleList {
		tabh := &proto.TimeAttackDungeonBattleHistoryDB{
			DungeonType:         proto.TimeAttackDungeonType(battleInfo.DungeonType),
			GeasId:              battleInfo.GeasId,
			DefaultPoint:        battleInfo.DefaultPoint,
			ClearTimePoint:      battleInfo.ClearTimePoint,
			EndFrame:            battleInfo.Frame,
			MainCharacterDBs:    make([]*proto.TimeAttackDungeonCharacterDB, 0),
			SupportCharacterDBs: make([]*proto.TimeAttackDungeonCharacterDB, 0),
		}
		for _, tadc := range battleInfo.MainCharacterDBs {
			tabh.MainCharacterDBs = append(tabh.MainCharacterDBs, GetTimeAttackDungeonCharacterDB(tadc))
		}
		for _, tadc := range battleInfo.SupportCharacterDBs {
			tabh.SupportCharacterDBs = append(tabh.SupportCharacterDBs, GetTimeAttackDungeonCharacterDB(tadc))
		}
		info.BattleHistoryDBs = append(info.BattleHistoryDBs, tabh)
	}
	for _, sweepTime := range bin.SweepTime {
		info.SweepHistoryDates = append(info.SweepHistoryDates, mx.Unix(sweepTime, 0))
	}
	return info
}

func GetTimeAttackDungeonCharacterDB(bin *sro.TimeAttackDungeonCharacter) *proto.TimeAttackDungeonCharacterDB {
	info := &proto.TimeAttackDungeonCharacterDB{
		ServerId:  bin.ServerId,
		UniqueId:  bin.UniqueId,
		CostumeId: bin.CostumeId,
		StarGrade: bin.StarGrade,
		Level:     bin.Level,
		HasWeapon: bin.HasWeapon,
		WeaponDB:  nil,
		IsAssist:  bin.IsAssist,
	}
	if bin.WeaponInfo != nil {
		info.WeaponDB = &proto.WeaponDB{
			Type:                   proto.ParcelType_CharacterWeapon,
			UniqueId:               bin.WeaponInfo.UniqueId,
			Level:                  bin.WeaponInfo.Level,
			Exp:                    bin.WeaponInfo.Exp,
			StarGrade:              bin.WeaponInfo.StarGrade,
			BoundCharacterServerId: bin.WeaponInfo.CharacterServerId,
		}
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

func TimeAttackDungeonClose(s *enter.Session) ([]*ParcelResult, int64) {
	bin := GetTimeAttackBin(s)
	curBin := GetTimeAttackRoom(s, bin.GetCurRoom())
	if curBin == nil {
		return nil, 0
	}
	curBin.RewardTime = time.Now().Unix()
	bin.CurRoom = 0
	// 计算总得分
	score := int64(0)
	for _, battleInfo := range curBin.BattleList {
		score += battleInfo.ClearTimePoint
		score += battleInfo.DefaultPoint
	}
	curBin.Score = score
	if curBin.IsPractice {
		// delete(bin.TimeAttackRoomList, curBin.RoomId)
		return nil, score
	}
	if previous := GetTimeAttackRoom(s, bin.GetPreviousRoom()); previous == nil || previous.Score <= score {
		bin.PreviousRoom = curBin.RoomId
	}
	// 获取奖励
	list := GetTimeAttackDungeonParcelResultByScore(score, curBin.SeasonId)

	return list, score
}

func GetTimeAttackDungeonParcelResultByScore(score, seasonId int64) (list []*ParcelResult) {
	list = make([]*ParcelResult, 0)
	conf := gdconf.GetTimeAttackDungeonSeasonManageExcelById(seasonId)
	if conf == nil {
		return
	}
	rewardConf := gdconf.GetTimeAttackDungeonRewardExcelTable(conf.TimeAttackDungeonRewardId)
	if rewardConf == nil {
		return
	}
	for index, minPoint := range rewardConf.RewardMinPoint {
		if minPoint > score {
			break
		}
		prInfo := &ParcelResult{
			ParcelType: proto.ParcelType_None.Value(rewardConf.RewardParcelType[index]),
			ParcelId:   rewardConf.RewardParcelId[index],
			Amount:     rewardConf.RewardParcelDefaultAmount[index],
		}
		if rewardConf.RewardType[index] == "TimeWeight" {
			prInfo.Amount = score
		}
		list = append(list, prInfo)
	}
	return
}
