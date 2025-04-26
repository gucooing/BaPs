package db

import (
	"errors"
	"fmt"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/driver/mysql"
	gromlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"

	"github.com/gucooing/BaPs/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type YostarRank struct {
	SeasonId int64
	Uid      int64   `gorm:"unique"`
	Score    float64 // 分数
	Rank     int64
}

func RaidUserTable(x int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(fmt.Sprintf("raid_rank_%v", x))
	}
}

func RaidEliminateUserTable(x int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(fmt.Sprintf("raid_eliminate_rank_%v", x))
	}
}

func ArenaUserTable(x int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(fmt.Sprintf("arena_rank_%v", x))
	}
}

func NewYostarRank(cfg *config.DB) *gorm.DB {
	switch cfg.DbType {
	case "sqlite":
		return NewSqlite(cfg.Dsn)
	case "mysql":
		return NewMysql(cfg.Dsn)
	default:
		log.Panicln("数据库的类型只支持 'sqlite' 和 'mysql' ")
	}

	return nil
}

// UpAllYostarRank 批量覆盖保存排名数据
func UpAllYostarRank(yostarRankSql *gorm.DB, x []*YostarRank, table func(tx *gorm.DB) *gorm.DB) error {
	err := yostarRankSql.Clauses(clause.OnConflict{
		UpdateAll: true, // 如果主键冲突,则更新所有字段
	}).Scopes(table).CreateInBatches(&x, 200).Error
	return err
}

// GetAllYostarRank 拉取全部排名数据
func GetAllYostarRank(yostarRankSql *gorm.DB, table func(tx *gorm.DB) *gorm.DB) []*YostarRank {
	r := make([]*YostarRank, 0)
	yostarRankSql.Scopes(table).Model(&YostarRank{}).Find(&r)
	return r
}

func NewMysql(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gromlogger.Default.LogMode(gromlogger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(errors.New("连接mysql数据库失败,请检查config中的配置和数据库是否存在"))
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(1000)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100000)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(100 * time.Millisecond) // 0.1 秒

	return db
}

func NewSqlite(dsn string) *gorm.DB {
	db, err := gorm.Open(gormlite.Open(dsn), &gorm.Config{
		Logger: gromlogger.Default.LogMode(gromlogger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(errors.New("连接sqlite数据库失败,请检查config中的配置和数据库目录是否存在"))
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(1000)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100000)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(100 * time.Millisecond) // 0.1 秒

	return db
}
