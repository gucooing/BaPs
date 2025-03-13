package sdk

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
)

type YostarCreateloginRequest struct {
	YostarToken    string `form:"yostar_token"`
	DeviceId       string `form:"deviceId"`
	CreateNew      int32  `form:"createNew"`
	YostarUsername string `form:"yostar_username"`
	YostarUid      int64  `form:"yostar_uid"`
	ChannelId      string `form:"channelId"`
}

type YostarCreateloginResponse struct {
	Result         int32  `json:"result"`
	Uid            string `json:"uid"`
	Token          string `json:"token"`
	YostarUid      string `json:"yostar_uid"`
	YostarUsername string `json:"yostar_username"`
	IsNew          int32  `json:"isNew"`
}

// YostarCreatelogin 登录完成验证
func (s *SDK) YostarCreatelogin(c *gin.Context) {
	req := &YostarCreateloginRequest{}
	rsp := &YostarCreateloginResponse{
		Result: -1,
	}
	defer c.JSON(200, rsp)
	err := c.ShouldBind(req)
	if err != nil {
		return
	}
	// 拉取YostarAccount
	yostarAccount := db.GetYostarAccountByYostarUid(req.YostarUid)
	if yostarAccount == nil {
		logger.Debug("拉取YostarAccount数据库信息失败")
		return
	}
	// 验证token
	if yostarAccount.YostarToken != req.YostarToken ||
		yostarAccount.YostarAccount != req.YostarUsername {
		logger.Error("邮箱:%s,YostarToken验证失败 user", req.YostarUsername)
		rsp.Result = 100140
		return
	}
	// 拉取YostarUser
	yostarUser := db.GetYostarUserByYostarUid(req.YostarUid)
	if yostarUser == nil {
		if !config.GetAutoRegistration() {
			logger.Debug("邮箱:%s,账号不存在且关闭自动注册  user", req.YostarUsername)
			return
		}
		logger.Debug("邮箱:%s,账号不存在进行自动注册 user", req.YostarUsername)
		yostarUser, err = db.AddYostarUserByYostarUid(req.YostarUid)
		if err != nil {
			logger.Debug("自动注册sdk账号失败,数据库错误:%s user", err.Error())
			return
		}
	}
	if yostarUser == nil {
		logger.Debug("邮箱:%s,进行数据库操作时候有未知错误 user", req.YostarUsername)
		return
	}
	// 更换设备后刷新token
	if yostarUser.DeviceId != req.DeviceId {
		yostarUser.Token = alg.RandStr(30)
	}
	yostarUser.ChannelId = req.ChannelId
	yostarUser.DeviceId = req.DeviceId
	if err = db.UpdateYostarUser(yostarUser); err != nil {
		logger.Debug("数据库写入出现错误:%s user", err.Error())
		return
	}
	rsp.Result = 0
	rsp.Token = yostarUser.Token
	rsp.Uid = strconv.Itoa(int(yostarUser.Uid))
	rsp.YostarUid = strconv.Itoa(int(yostarUser.YostarUid))
	rsp.YostarUsername = req.YostarUsername
	logger.Debug("邮箱:%s,登录成功 YostarUid:%v,Token:%s,Uid:%v", req.YostarUsername, req.YostarUid, yostarUser.Token, yostarUser.Uid)
}

type YostarLoginRequest struct {
	DeviceId      string `form:"deviceId"`
	Uid           int64  `form:"uid"`
	StoreId       string `form:"storeId"`
	Platform      string `form:"platform"`
	CaptchaOutput string `form:"captcha_output"`
	GenTime       int64  `form:"gen_time"`
	Token         string `form:"token"`
	CaptchaId     string `form:"captcha_id"`
	LotNumber     string `form:"lot_number"`
	PassToken     string `form:"pass_token"`
}

type YostarLoginResponse struct {
	Result             int32       `json:"result"`
	AccessToken        string      `json:"accessToken"`
	Birth              interface{} `json:"birth"`
	YostarUid          string      `json:"yostar_uid"`
	YostarUsername     string      `json:"yostar_username"`
	Transcode          string      `json:"transcode"`
	CurrentTimestampMs int64       `json:"current_timestamp_ms"`
	Check7Until        int32       `json:"check7until"`
	Migrated           bool        `json:"migrated"`
	ShowMigratePage    bool        `json:"show_migrate_page"`
	ChannelId          string      `json:"channelId"`
	KrKmcStatus        int32       `json:"kr_kmc_status"`
}

// YostarLogin 获取登录网关的token
func (s *SDK) YostarLogin(c *gin.Context) {
	req := &YostarLoginRequest{}
	rsp := &YostarLoginResponse{
		Result: -1,
	}
	defer c.JSON(200, rsp)
	err := c.ShouldBind(req)
	if err != nil {
		return
	}
	// 拉取YostarUser
	yostarUser := db.GetYostarUserByUid(req.Uid)
	if yostarUser == nil {
		logger.Debug("UID:%v,未知的登录请求", req.Uid)
		rsp.Result = 1
		return
	}
	// Token验证
	if yostarUser.Token != req.Token {
		logger.Debug("UID:%v,ToKen:%s,Token验证失败", req.Uid, req.Token)
		rsp.Result = 1
		return
	}
	// 拉取YostarAccount
	yostarAccount := db.GetYostarAccountByYostarUid(yostarUser.YostarUid)
	if yostarAccount == nil {
		logger.Debug("拉取YostarAccount数据库信息失败 login")
		return
	}
	// 拉取YoStarUserLogin
	yoStarUserLogin := db.GetYoStarUserLoginByYostarUid(yostarAccount.YostarUid)
	if yoStarUserLogin == nil {
		if !config.GetAutoRegistration() {
			logger.Debug("邮箱:%s,账号不存在且关闭自动注册  login", yostarAccount.YostarAccount)
			return
		}
		yoStarUserLogin, err = db.AddYoStarUserLoginByYostarUid(yostarAccount.YostarUid)
		if err != nil {
			logger.Debug("自动注册登录账号失败,数据库错误:%s login", err.Error())
			return
		}
	}
	if yoStarUserLogin == nil {
		logger.Debug("邮箱:%s,进行数据库操作时候有未知错误 login", yostarAccount.YostarAccount)
		return
	}
	// 黑名单验证
	if yoStarUserLogin.Ban {
		logger.Debug("邮箱:%s,账号已被封禁,原因:%s", yostarAccount.YostarAccount, yoStarUserLogin.BanMsg)
		rsp.Result = 100305
		return
	}
	// 设备黑名单
	if blackDevice := db.GetBlackDeviceByYostarUid(req.DeviceId); blackDevice != nil ||
		req.DeviceId == "" {
		logger.Debug("邮箱:%s,DeviceId:%s,设备已被封禁", yostarAccount.YostarAccount, req.DeviceId)
		rsp.Result = 100100
		return
	}
	yoStarUserLogin.YostarLoginToken = alg.RandStr(30)
	// 更新YoStarUserLogin
	if err = db.UpdateYoStarUserLogin(yoStarUserLogin); err != nil {
		logger.Debug("数据库写入出现错误:%s login", err.Error())
		return
	}
	// 拉取游戏数据
	rsp.Result = 0
	rsp.ChannelId = yostarUser.ChannelId
	rsp.CurrentTimestampMs = time.Now().UnixMicro()
	rsp.KrKmcStatus = 2
	rsp.Migrated = true
	rsp.YostarUid = strconv.Itoa(int(yostarUser.YostarUid))
	rsp.YostarUsername = yostarAccount.YostarAccount
	rsp.AccessToken = yoStarUserLogin.YostarLoginToken
}
