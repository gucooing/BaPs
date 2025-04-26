package db

import (
	"github.com/gucooing/BaPs/config"
	dbgorm "github.com/gucooing/BaPs/db/gorm"
	dbstruct "github.com/gucooing/BaPs/db/struct"
	"github.com/gucooing/BaPs/pkg/logger"
	_ "github.com/ncruces/go-sqlite3/embed"
	"log"
)

// DBGame 接口
type DBGame interface {
	// BlackDevice
	GetBlackDeviceByYostarUid(deviceId string) *dbstruct.BlackDevice // 使用DeviceId拉取数据
	// YostarAccount
	GetYostarAccountByYostarAccount(yostarAccount string) *dbstruct.YostarAccount          // 使用YostarAccount拉取数据
	GetYostarAccountByYostarUid(yostarUid int64) *dbstruct.YostarAccount                   // 使用YostarUid拉取数据
	AddYostarAccountByYostarAccount(yostarAccount string) (*dbstruct.YostarAccount, error) // 指定YostarAccount创建数据
	UpdateYostarAccount(data *dbstruct.YostarAccount) error                                // 更新账号数据
	// YostarClan
	GetYostarClanByServerId(serverId int64) *dbstruct.YostarClan           // 使用ServerId拉取数据
	GetYostarClanByClanName(clanName string) *dbstruct.YostarClan          // 使用ClanName拉取数据
	AddYostarClanByClanName(clanName string) (*dbstruct.YostarClan, error) // 指定ClanName创建数据
	UpdateYostarClan(data *dbstruct.YostarClan) error                      // 更新社团数据
	UpAllYostarClan(list []*dbstruct.YostarClan) error                     // 批量覆盖保存社团数据
	// YostarFriend
	GetYostarFriendByAccountServerId(accountServerId int64) *dbstruct.YostarFriend    // 使用AccountServerId拉取数据
	AddYostarFriendByYostarUid(accountServerId int64) (*dbstruct.YostarFriend, error) // 指定YostarUid创建数据
	UpdateYostarFriend(data *dbstruct.YostarFriend) error                             // 更新好友数据
	UpAllYostarFriend(list []*dbstruct.YostarFriend) error                            // 批量覆盖保存好友数据
	// YostarGame
	GetYostarGameByAccountServerId(accountServerId int64) *dbstruct.YostarGame    // 使用AccountServerId拉取数据
	AddYostarGameByYostarUid(accountServerId int64) (*dbstruct.YostarGame, error) // 指定YostarUid创建数据
	UpdateYostarGame(data *dbstruct.YostarGame) error                             // 更新账号数据
	UpAllYostarGame(list []*dbstruct.YostarGame) error                            // 批量覆盖保存账号数据
	// YostarMail
	GetAllYostarMail() []*dbstruct.YostarMail                          //拉取全部邮件
	AddYostarMailBySender(sender string) (*dbstruct.YostarMail, error) // 指定Sender创建数据
	UpdateYostarMail(data *dbstruct.YostarMail) error                  // 更新邮件数据
	// YostarUser
	GetYostarUserByUid(uid int64) *dbstruct.YostarUser                      // 使用Uid拉取数据
	GetYostarUserByYostarUid(yostarUid int64) *dbstruct.YostarUser          // 使用YostarUid拉取数据
	AddYostarUserByYostarUid(yostarUid int64) (*dbstruct.YostarUser, error) // 指定YostarAccount创建数据
	UpdateYostarUser(data *dbstruct.YostarUser) error                       // 更新账号数据
	// YostarUserLogin
	GetYoStarUserLoginByYostarUid(yostarUid int64) *dbstruct.YostarUserLogin          // 使用YostarUid拉取数据
	AddYoStarUserLoginByYostarUid(yostarUid int64) (*dbstruct.YostarUserLogin, error) // 指定YostarUid创建数据
	UpdateYoStarUserLogin(data *dbstruct.YostarUserLogin) error                       // 更新账号数据
}

var sqlGame *DBGameService

type DBGameService struct {
	db DBGame
}

func NewDBService(cfg *config.DB) {
	d := new(DBGameService)
	var err error
	switch cfg.DbType {
	case "sqlite", "mysql":
		d.db, err = dbgorm.NewDbGorm(cfg.DbType, cfg.Dsn)
	default:
		log.Panicln("数据库的类型只支持 'sqlite' 和 'mysql' ")
		return
	}
	if err != nil {
		log.Panicln(err)
		return
	}
	sqlGame = d
	logger.Info("数据库连接成功")
}

func GetDBGame() DBGame {
	return sqlGame.db
}
