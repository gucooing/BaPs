package config

import (
	"encoding/json"
	"errors"
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
}

type GateWay struct {
	MaxPlayerNum   int64           `json:"MaxPlayerNum"`
	BlackCmd       map[string]bool `json:"BlackCmd"`
	IsLogMsgPlayer bool            `json:"IsLogMsgPlayer"`
	IsToken        bool            `json:"IsToken"`
	GetTokenUrl    string          `json:"GetTokenUrl"`
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

var CONF *Config = nil

func SetDefaultConfig() {
	CONF = DefaultConfig
}

func GetConfig() *Config {
	return CONF
}

func GetGucooingApiKey() string {
	return CONF.GucooingApiKey
}

func GetAutoRegistration() bool {
	return CONF.AutoRegistration
}

func GetHttpNet() *HttpNet {
	return CONF.HttpNet
}

func GetGateWay() *GateWay {
	return GetConfig().GateWay
}

var FileNotExist = errors.New("config file not found")

func LoadConfig() error {
	filePath := "./config.json"
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
		CertFile:  "./data/cert.pem",
		KeyFile:   "./data/key.pem",
	},
	GateWay: &GateWay{
		MaxPlayerNum:   0,
		BlackCmd:       make(map[string]bool),
		IsLogMsgPlayer: false,
		IsToken:        true,
		GetTokenUrl:    "http://127.0.0.1:8080/gucooing/api/getToken/ba",
	},
	DB: &DB{
		DbType: "sqlite",
		Dsn:    "BaPs.db",
	},
}
