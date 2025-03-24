package enter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
	pb "google.golang.org/protobuf/proto"
)

var MaxCachePlayerTime = 120 // 最大玩家缓存时间 单位:分钟
var MaxPlayerNum int64 = 0   // 最大在线玩家

type Session struct {
	AccountServerId int64
	YostarUID       int64
	MxToken         string
	EndTime         time.Time
	AccountState    proto.AccountState
	PlayerBin       *sro.PlayerBin // 玩家数据
	Actions         map[proto.ServerNotificationFlag]bool
	GoroutinesSync  sync.Mutex
	AccountFriend   *AccountFriend
	Mission         *Mission
	Toast           []string
	PlayerHash      map[int64]any
	arenaInfo       *ArenaInfo // 竞技场临时数据
}

// 定时检查一次是否有用户长时间离线
func (e *EnterSet) checkSession() {
	yostarGameList := make([]*db.YostarGame, 0)
	for accountServerId, info := range GetAllSession() {
		if time.Now().After(info.EndTime) {
			bin := info.GetYostarGame()
			if bin != nil {
				yostarGameList = append(yostarGameList, bin)
			}
			DelSession(accountServerId)
			logger.Debug("AccountId:%v,超时离线", accountServerId)
		}
	}
	if len(yostarGameList) == 0 {
		return
	}
	if db.UpAllYostarGame(yostarGameList) != nil {
		logger.Error("玩家数据保存失败")
	} else {
		logger.Info("玩家数据保存完毕")
	}
}

// GetSessionBySessionKey 获取指定在线玩家
func GetSessionBySessionKey(sessionKey *proto.SessionKey) *Session {
	if sessionKey == nil ||
		sessionKey.AccountServerId == 0 ||
		sessionKey.MxToken == "" {
		return nil
	}
	e := getEnterSet()
	e.sessionSync.RLock()
	defer e.sessionSync.RUnlock()
	if info, ok := e.SessionMap[sessionKey.AccountServerId]; ok {
		if info.MxToken != sessionKey.MxToken {
			return nil
		}
		return info
	}
	return nil
}

// GetSessionByAccountServerId 获取指定在线玩家
func GetSessionByAccountServerId(accountServerId int64) *Session {
	if accountServerId == 0 {
		return nil
	}
	e := getEnterSet()
	e.sessionSync.RLock()
	defer e.sessionSync.RUnlock()
	if info, ok := e.SessionMap[accountServerId]; ok {
		return info
	}
	return nil
}

// GetSessionByUid 此接口的用处是拉取玩家数据,包括在数据库中的
func GetSessionByUid(uid int64) *Session {
	if uid == 0 {
		return nil
	}
	e := getEnterSet()
	e.sessionSync.RLock()
	info, ok := e.SessionMap[uid]
	if ok {
		e.sessionSync.RUnlock()
		return info
	}
	e.sessionSync.RUnlock()
	bin := db.GetYostarGameByAccountServerId(uid)
	if bin == nil || bin.BinData == nil {
		return nil
	}
	info = NewSession(uid)
	info.AccountServerId = uid
	info.EndTime = time.Now().Add(time.Duration(MaxCachePlayerTime) * time.Minute)
	err := pb.Unmarshal(bin.BinData, info.PlayerBin)
	if err != nil || info.PlayerBin.GetBaseBin().GetAccountId() != uid {
		return nil
	}
	AddSession(info)
	return info
}

func NewSession(accountServerId int64) *Session {
	return &Session{
		AccountServerId: accountServerId,
		Actions:         make(map[proto.ServerNotificationFlag]bool),
		PlayerBin:       new(sro.PlayerBin),
		GoroutinesSync:  sync.Mutex{},
		AccountFriend:   GetAccountFriend(accountServerId),
	}
}

// GetAllSession 获取全部在线玩家
func GetAllSession() map[int64]*Session {
	allSession := make(map[int64]*Session)
	e := getEnterSet()
	e.sessionSync.RLock()
	defer e.sessionSync.RUnlock()
	for k, v := range e.SessionMap {
		allSession[k] = v
	}
	return allSession
}

// GetAllSessionList 获取全部在线玩家-列表
func GetAllSessionList() []*Session {
	allSession := make([]*Session, 0)
	e := getEnterSet()
	e.sessionSync.RLock()
	defer e.sessionSync.RUnlock()
	for _, v := range e.SessionMap {
		allSession = append(allSession, v)
	}
	return allSession
}

// GetSessionNum 获取在线玩家数量
func GetSessionNum() int64 {
	e := getEnterSet()
	e.sessionSync.RLock()
	defer e.sessionSync.RUnlock()
	return int64(len(e.SessionMap))
}

// DelSession 删除指定在线玩家
func DelSession(accountServerId int64) bool {
	e := getEnterSet()
	e.sessionSync.Lock()
	defer e.sessionSync.Unlock()
	if e.SessionMap == nil {
		e.SessionMap = make(map[int64]*Session)
	}
	if _, ok := e.SessionMap[accountServerId]; ok {
		delete(e.SessionMap, accountServerId)
		return true
	}
	return false
}

// AddSession 添加Session
func AddSession(x *Session) bool {
	if x == nil ||
		x.AccountServerId == 0 {
		return false
	}
	e := getEnterSet()
	e.sessionSync.Lock()
	defer e.sessionSync.Unlock()
	if e.SessionMap == nil {
		e.SessionMap = make(map[int64]*Session)
	}
	if _, ok := e.SessionMap[x.AccountServerId]; ok {
		e.SessionMap[x.AccountServerId] = x
		return false
	}
	e.SessionMap[x.AccountServerId] = x
	return true
}

// Close 保存全部玩家数据
func Close() {
	// 保存玩家主要数据
	yostarGameList := make([]*db.YostarGame, 0)
	for _, info := range GetAllSession() {
		bin := info.GetYostarGame()
		if bin != nil {
			yostarGameList = append(yostarGameList, bin)
		}
	}
	if db.UpAllYostarGame(yostarGameList) != nil {
		logger.Error("玩家数据保存失败")
	} else {
		logger.Info("玩家数据保存完毕")
	}
	// 保存玩家次要数据 (好友数据
	yostarFriendList := make([]*db.YostarFriend, 0)
	for _, info := range GetAllAccountFriend() {
		bin := info.GetYostarFriend()
		if bin != nil {
			yostarFriendList = append(yostarFriendList, bin)
		}
	}
	if db.UpAllYostarFriend(yostarFriendList) != nil {
		logger.Error("好友数据保存失败")
	} else {
		logger.Info("好友数据保存完毕")
	}
	// 保存社团数据
	yostarClanList := make([]*db.YostarClan, 0)
	for _, info := range GetAllYostarClan() {
		bin := info.GetYostarClan()
		if bin != nil {
			yostarClanList = append(yostarClanList, bin)
		}
	}
	if db.UpAllYostarClan(yostarClanList) != nil {
		logger.Error("社团数据保存失败")
	} else {
		logger.Info("社团数据保存完毕")
	}
}

// GetPbBinData 将玩家pb数据转二进制数据
func (x *Session) GetPbBinData() []byte {
	bin, err := pb.Marshal(x.PlayerBin)
	if err != nil {
		return nil
	}
	return bin
}

// GetYostarGame 将玩家数据转成数据库格式
func (x *Session) GetYostarGame() *db.YostarGame {
	if x == nil {
		return nil
	}
	x.GoroutinesSync.Lock() // 唯一线程操作锁
	var fin = true
	defer func() {
		x.GoroutinesSync.Unlock()
		if !fin {
			logger.Debug("玩家:%v,数据保存失败,数据保存将保存到服务端硬盘,下次启动时将尝试写入数据库", x.AccountServerId)
			if err := x.upDataDisk(); err != nil {
				logger.Debug("玩家:%v,数据保存将到服务端硬盘失败,失败原因:%s", x.AccountServerId, err.Error())
			}
		}
	}()
	bin, err := pb.Marshal(x.PlayerBin)
	if err != nil {
		fin = false
		return nil
	}
	data := &db.YostarGame{
		AccountServerId: x.AccountServerId,
		BinData:         bin,
	}

	return data
}

// 将玩家二进制数据写入磁盘中
func (x *Session) upDataDisk() error {
	bin, err := pb.Marshal(x.PlayerBin)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("./player/%v.bin", x.AccountServerId), bin, 0644)
	return err
}

// TaskUpDiskPlayerData 将磁盘中的二进制玩家数据写入数据库中
func TaskUpDiskPlayerData() bool {
	files, err := filepath.Glob(filepath.Join("./player/", "*.bin"))
	if err != nil {
		return true
	}
	logger.Info("尝试保存本地玩家数据")
	for _, file := range files {
		bin, err := os.ReadFile(file)
		if err != nil {
			logger.Error("读取本地玩家数据失败:%s", err.Error())
			return false
		}
		accountServerId := alg.S2I64(strings.TrimSuffix(filepath.Base(file), ".bin"))
		if accountServerId == 0 {
			logger.Error("本地玩家数据文件名错误,文件:%s", file)
			return false
		}
		data := &db.YostarGame{
			AccountServerId: accountServerId,
			BinData:         bin,
		}
		if err = db.UpdateYostarGame(data); err != nil {
			return false
		}
		if os.Remove(file) != nil {
			logger.Warn("删除本地玩家数据文件失败,文件:%s,可能是权限不足导致的,请手动删除,避免下次启动时数据被覆盖", file)
		}
	}
	logger.Info("保存本地玩家数据成功")
	return true
}

func (x *Session) getPlayerHash() map[int64]any {
	if x.PlayerHash == nil {
		x.PlayerHash = make(map[int64]any)
	}
	return x.PlayerHash
}

// AddPlayerHash 写入数据到哈希表中
func (x *Session) AddPlayerHash(k int64, v any) bool {
	list := x.getPlayerHash()
	if _, ok := list[k]; ok {
		return false
	}
	list[k] = v
	return true
}

func (x *Session) getPlayerHashByKeyId(k int64) any {
	list := x.getPlayerHash()
	return list[k]
}

func (x *Session) GetCharacterByKeyId(k int64) *sro.CharacterInfo {
	v := x.getPlayerHashByKeyId(k)
	switch info := v.(type) {
	case *sro.CharacterInfo:
		return info
	default:
		return nil
	}
}

func (x *Session) GetItemByKeyId(k int64) *sro.ItemInfo {
	v := x.getPlayerHashByKeyId(k)
	switch info := v.(type) {
	case *sro.ItemInfo:
		return info
	default:
		return nil
	}
}

func (x *Session) GetEquipmentByKeyId(k int64) *sro.EquipmentInfo {
	v := x.getPlayerHashByKeyId(k)
	switch info := v.(type) {
	case *sro.EquipmentInfo:
		return info
	default:
		return nil
	}
}
