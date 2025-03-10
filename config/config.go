package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	LogLevel         string   `json:"LogLevel"`
	ResourcesPath    string   `json:"ResourcesPath"`
	DataPath         string   `json:"DataPath"`
	GucooingApiKey   string   `json:"GucooingApiKey"`
	AutoRegistration bool     `json:"AutoRegistration"`
	HttpNet          *HttpNet `json:"HttpNet"`
	GateWay          *GateWay `json:"GateWay"`
	DB               *DB      `json:"DB"`
	Irc              *Irc     `json:"Irc"`
	RankDB           *DB      `json:"RankDB"`
}

type GateWay struct {
	MaxPlayerNum       int64           `json:"MaxPlayerNum"`
	MaxCachePlayerTime int             `json:"MaxCachePlayerTime"`
	BlackCmd           map[string]bool `json:"BlackCmd"`
	IsLogMsgPlayer     bool            `json:"IsLogMsgPlayer"`
}

type HttpNet struct {
	InnerAddr string `json:"InnerAddr"`
	InnerPort string `json:"InnerPort"`
	OuterAddr string `json:"OuterAddr"`
	OuterPort string `json:"OuterPort"`
	Tls       bool   `json:"Tls"`
	CertFile  string `json:"CertFile"`
	KeyFile   string `json:"KeyFile"`
}

type DB struct {
	DbType string `json:"dbType"`
	Dsn    string `json:"dsn"`
}

type Irc struct {
	HostAddress string `json:"HostAddress"`
	Port        int32  `json:"Port"`
	Password    string `json:"Password"`
}

var CONF *Config = nil

func SetDefaultConfig() {
	log.Printf("config不存在,使用默认配置\n")
	CONF = DefaultConfig
}

func GetConfig() *Config {
	if CONF == nil {
		SetDefaultConfig()
	}
	return CONF
}

func GetGucooingApiKey() string {
	return GetConfig().GucooingApiKey
}

func GetAutoRegistration() bool {
	return GetConfig().AutoRegistration
}

func GetHttpNet() *HttpNet {
	return GetConfig().HttpNet
}

func GetGateWay() *GateWay {
	return GetConfig().GateWay
}

func GetIsLogMsgPlayer() bool {
	return GetConfig().GateWay.IsLogMsgPlayer
}

func GetBlackCmd() map[string]bool {
	return GetConfig().GateWay.BlackCmd
}

func GetRankDB() *DB {
	return GetConfig().RankDB
}

func GetIrc() *Irc {
	return GetConfig().Irc
}

var FileNotExist = errors.New("config file not found")

func LoadConfig() error {
	filePath := "./config/config.json"
	f, err := os.Open(filePath)
	if err != nil {
		return FileNotExist
	}
	defer func() {
		_ = f.Close()
	}()
	c := new(Config)
	d := json.NewDecoder(f)
	if err := d.Decode(c); err != nil {
		return err
	}
	CONF = c
	return nil
}

var DefaultConfig = &Config{
	LogLevel:         "Info",
	ResourcesPath:    "./resources",
	DataPath:         "./data",
	GucooingApiKey:   "123456",
	AutoRegistration: true,
	HttpNet: &HttpNet{
		InnerAddr: "0.0.0.0",
		InnerPort: "5000",
		OuterAddr: "127.0.0.1",
		OuterPort: "5000",
		Tls:       false,
		CertFile:  "./config/cert.pem",
		KeyFile:   "./config/key.pem",
	},
	GateWay: &GateWay{
		MaxPlayerNum:       0,
		MaxCachePlayerTime: 720,
		BlackCmd:           make(map[string]bool),
		IsLogMsgPlayer:     false,
	},
	DB: &DB{
		DbType: "sqlite",
		Dsn:    "./config/BaPs.db",
	},
	RankDB: &DB{
		DbType: "sqlite",
		Dsn:    "./config/Rank.db",
	},
	Irc: &Irc{
		HostAddress: "127.0.0.1",
		Port:        16666,
		Password:    "mx123",
	},
}
