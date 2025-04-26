package db_gorm

import (
	"errors"
	dbstruct "github.com/gucooing/BaPs/db/struct"
	"time"

	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gromlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DbGorm struct {
	sql *gorm.DB
}

func NewDbGorm(dbType, dsn string) (*DbGorm, error) {
	x := new(DbGorm)
	var err error
	switch dbType {
	case "sqlite":
		err = x.NewSqlite(dsn)
	case "mysql":
		err = x.NewMysql(dsn)
	default:
		err = errors.New("uknown DB type")
	}

	if err != nil {
		return nil, err
	}

	err = x.sql.AutoMigrate(
		&dbstruct.YostarAccount{},
		&dbstruct.BlackDevice{},
		&dbstruct.YostarClan{},
		&dbstruct.YostarFriend{},
		&dbstruct.YostarGame{},
		&dbstruct.YostarMail{},
		&dbstruct.YostarUser{},
		&dbstruct.YostarUserLogin{},
	)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *DbGorm) NewMysql(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gromlogger.Default.LogMode(gromlogger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(1000)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100000)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(100 * time.Millisecond) // 0.1 秒
	x.sql = db

	return nil
}

func (x *DbGorm) NewSqlite(dsn string) error {
	db, err := gorm.Open(gormlite.Open(dsn), &gorm.Config{
		Logger: gromlogger.Default.LogMode(gromlogger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(1000)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100000)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(100 * time.Millisecond) // 0.1 秒
	x.sql = db

	return nil
}
