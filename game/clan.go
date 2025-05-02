package game

import (
	"math/rand"

	"github.com/gucooing/BaPs/common/check"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetClanServerId(s *enter.Session) int64 {
	bin := GetBaseBin(s)
	if bin == nil {
		return 0
	}
	return bin.ClanId
}

func SetClanServerId(s *enter.Session, clanId int64) {
	bin := GetBaseBin(s)
	if bin == nil {
		return
	}
	bin.ClanId = clanId
}

func GetIrcServerConfig(s *enter.Session) *proto.IrcServerConfig {
	cfg := config.GetIrc()
	if cfg == nil {
		return nil
	}
	return &proto.IrcServerConfig{
		HostAddress: cfg.HostAddress,
		Port:        cfg.Port,
		Password:    cfg.Password,
	}
}

func GetClanDB(clanInfo *enter.YostarClan) *proto.ClanDB {
	if clanInfo == nil {
		return nil
	}
	info := &proto.ClanDB{
		ClanDBId:        clanInfo.ServerId,
		ClanName:        clanInfo.ClanName,
		ClanChannelName: clanInfo.ClanName,
		ClanNotice:      clanInfo.Notice,
		ClanMemberCount: clanInfo.GetMemberCount(),
		ClanJoinOption:  proto.ClanJoinOption(clanInfo.JoinOption),
	}
	ps := enter.GetSessionByUid(clanInfo.President)
	if ps == nil {
		return nil
	}
	info.ClanPresidentNickName = GetNickname(ps)
	info.ClanPresidentRepresentCharacterUniqueId = GetRepresentCharacterUniqueId(ps)
	info.ClanPresidentRepresentCharacterCostumeId = 0
	return info
}

func GetDefaultExposedClanDBs(s *enter.Session) []*proto.ClanDB {
	myClanInfo := enter.GetYostarClanByServerId(GetClanServerId(s))
	if myClanInfo != nil {
		return nil
	}
	allClanInfo := enter.GetAllYostarClanList()
	maxNum := alg.MainInt(30, len(allClanInfo)-1)
	clanList := make(map[int64]*enter.YostarClan, 0)
	for i := 0; i < maxNum; i++ {
		clanInfo := allClanInfo[rand.Intn(len(allClanInfo))]
		if clanInfo == nil ||
			clanList[clanInfo.ServerId] != nil {
			continue
		}
		clanList[clanInfo.ServerId] = clanInfo
	}
	list := make([]*proto.ClanDB, 0)
	for _, clan := range allClanInfo {
		list = append(list, GetClanDB(clan))
	}
	return list
}

func GetClanMemberDBs(clanInfo *enter.YostarClan) []*proto.ClanMemberDB {
	list := make([]*proto.ClanMemberDB, 0)
	for uid := range clanInfo.GetAllAccount() {
		ps := enter.GetSessionByUid(uid)
		list = append(list, GetClanMemberDB(ps))
	}

	return list
}

func GetClanMemberDB(s *enter.Session) *proto.ClanMemberDB {
	bin := GetBaseBin(s)
	if bin == nil {
		return nil
	}
	clanInfo := enter.GetYostarClanByServerId(bin.GetClanId())
	if clanInfo == nil {
		return &proto.ClanMemberDB{
			AccountId: bin.GetAccountId(),
		}
	}
	info := &proto.ClanMemberDB{
		ClanDBId:                   bin.GetClanId(),
		AccountId:                  bin.GetAccountId(),
		AccountLevel:               int64(GetAccountLevel(s)),
		AccountNickName:            GetNickname(s),
		AttachmentDB:               GetAccountAttachmentDB(s),
		CafeComfortValue:           9000,
		RepresentCharacterUniqueId: GetRepresentCharacterUniqueId(s),
		GameLoginDate:              GetLastConnectTime(s),
	}

	ca := clanInfo.GetClanAccount(bin.GetAccountId())
	info.AttendanceCount = ca.GetAttendanceCount()
	info.ClanSocialGrade = proto.ClanSocialGrade(ca.GetSocialGrade())
	info.JoinDate = ca.GetJoinDate()
	info.LastLoginDate = ca.GetLastLoginTime()
	return info
}

func NewClan(s *enter.Session, clanName string, joinOption proto.ClanJoinOption) proto.WebAPIErrorCode {
	bin := GetBaseBin(s)
	if bin == nil ||
		enter.GetYostarClanByClanName(clanName) != nil {
		logger.Debug("社团重复")
		return 15022
	}
	if !check.CheckName(clanName) {
		logger.Debug("社团名称检查不通过")
		return proto.WebAPIErrorCode_ClanNameWithInvalidLength
	}
	_, err := db.GetDBGame().AddYostarClanByClanName(clanName)
	if err != nil {
		logger.Debug("新社团数据库写入失败")
		return 15022
	}
	clanInfo := enter.GetYostarClanByClanName(clanName)
	if clanInfo == nil {
		logger.Debug("拉取新社团失败")
		return proto.WebAPIErrorCode_ClanNotFound
	}
	// 添加社长
	if !clanInfo.AddAccount(s.AccountServerId, int32(proto.ClanSocialGrade_President)) {
		logger.Debug("新社团设置社长失败")
		return 15022
	}
	// 设置社长
	if !clanInfo.SetPresident(s.AccountServerId) {
		logger.Debug("新社团设置社长失败")
		return 15022
	}
	clanInfo.SetNotice("欢迎游玩BaPs,这是一个半开源的免费服务器\n----By gucooing")
	SetClanServerId(s, clanInfo.ServerId)
	return 0
}

func SetLastLoginTime(s *enter.Session) {
	clanInfo := enter.GetYostarClanByServerId(GetClanServerId(s))
	if clanInfo == nil {
		return
	}
	ca := clanInfo.GetClanAccount(s.AccountServerId)
	if ca == nil {
		return
	}
	ca.SetLastLoginTime()
}

func GetClanName(s *enter.Session) string {
	clanInfo := enter.GetYostarClanByServerId(GetClanServerId(s))
	if clanInfo == nil {
		return ""
	}
	return clanInfo.ClanName
}
