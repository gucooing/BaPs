package enter

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	pb "google.golang.org/protobuf/proto"
)

type Session struct {
	AccountServerId int64
	MxToken         string
	EndTime         time.Time
	AccountState    proto.AccountState
	PlayerBin       *sro.PlayerBin // 玩家数据
}

// 定时检查一次是否有用户长时间离线
func (e *EnterSet) checkSession() {
	for accountServerId, info := range GetAllSession() {
		if time.Now().After(info.EndTime) {
			info.UpDate()
			DelSession(accountServerId)
			logger.Debug("AccountId:%v,超时离线", accountServerId)
		}
	}
}

// 有锁 检查并处理重复登录
func (e *EnterSet) checkSessionRepeat(accountServerId int64) {
	e.sessionSync.Lock()
	defer e.sessionSync.Unlock()
	if e.SessionMap == nil {
		e.SessionMap = make(map[int64]*Session)
	}
	if _, ok := e.SessionMap[accountServerId]; ok {
		delete(e.SessionMap, accountServerId)
	}
}

// GetSessionBySessionKey 获取指定在线玩家
func GetSessionBySessionKey(sessionKey *mx.SessionKey) *Session {
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

// UpAllDate 保存全部玩家数据
func UpAllDate() {
	for _, info := range GetAllSession() {
		info.UpDate()
	}
}

// UpDate 保存玩家数据
func (x *Session) UpDate() bool {
	var fin = true
	defer func() {
		if !fin {
			logger.Debug("玩家:%v,数据保存失败,数据保存将到服务端硬盘,将在下次启动时尝试写入数据库", x.AccountServerId)
			if err := x.upDataDisk(); err != nil {
				logger.Debug("玩家:%v,数据保存将到服务端硬盘失败,失败原因:%s", x.AccountServerId, err.Error())
			}
		}
	}()
	bin, err := pb.Marshal(x.PlayerBin)
	if err != nil {
		fin = false
		return false
	}
	data := &db.YostarGame{
		AccountServerId: x.AccountServerId,
		BinData:         bin,
	}
	if err = db.UpdateYostarGame(data); err != nil {
		fin = false
		return false
	}
	return true
}

func (x *Session) upDataDisk() error {
	bin, err := pb.Marshal(x.PlayerBin)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("./player/%v.bin", x.AccountServerId), bin, 0644)
	return err
}

// TaskUpDiskPlayerData 保存上次保存失败的数据到数据库中
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
		accountServerId := alg.S2I64(filepath.Base(file))
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
	}
	logger.Info("保存本地玩家数据成功")
	return true
}
