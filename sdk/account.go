package sdk

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
)

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type YostarAuthRequest struct {
	Account  string `form:"account"`
	AuthLang string `form:"authlang"`
	Platform string `form:"platform"`
}
type YostarAuthResponse struct {
	Result int32 `json:"result"`
}

// YostarAuthRequest 邮箱登录自动注册
func (s *SDK) YostarAuthRequest(c *gin.Context) {
	req := &YostarAuthRequest{}
	rsp := &YostarAuthResponse{
		Result: -1,
	}
	defer c.JSON(200, rsp)
	if err := c.ShouldBind(req); err != nil {
		return
	}
	// 验证邮箱是否符合规定//不需要了
	if !emailRe.MatchString(req.Account) {
		logger.Debug("邮箱:%s,格式错误", req.Account)
		rsp.Result = 100300
		return
	}
	logger.Debug("邮箱:%s,语言:%s,平台:%s 发起邮箱验证码登录", req.Account, req.AuthLang, req.Platform)
	// 验证请求间隔是否过短
	if code := s.code.GetCodeInfo(req.Account); code != nil {
		logger.Debug("邮箱:%s,验证码还未过期/使用不刷新", req.Account)
	} else {
		newCode := alg.RandCode()
		if err := s.code.SetCode(req.Account, newCode); err != nil {
			logger.Debug("邮箱:%s,验证码添加失败:%s", req.Account, err.Error())
			rsp.Result = 100302
			return
		} else {
			logger.Debug("邮箱:%s,写入新验证码:%v", req.Account, newCode)
		}
	}
	rsp.Result = 0
}

type YostarAuthSubmitRequest struct {
	Account string `form:"account"`
	Code    int32  `form:"code"`
}

type YostarAuthSubmitResponse struct {
	Result        int32  `json:"result"`
	YostarUid     string `json:"yostar_uid"`
	YostarToken   string `json:"yostar_token"`
	YostarAccount string `json:"yostar_account"`
}

// YostarAuthSubmit 邮箱登录验证验证码是否有效
func (s *SDK) YostarAuthSubmit(c *gin.Context) {
	req := &YostarAuthSubmitRequest{}
	rsp := &YostarAuthSubmitResponse{
		Result: -1,
	}
	defer c.JSON(200, rsp)
	err := c.ShouldBind(req)
	if err != nil {
		return
	}
	// 验证验证码是否有效
	if code := s.code.GetCodeInfo(req.Account); code == nil ||
		req.Code != code.Code {
		rsp.Result = 100306
		logger.Debug("邮箱:%s,无验证码", req.Account)
		return
	} else if req.Code != code.Code {
		rsp.Result = 100307
		code.FialNum++
		logger.Debug("邮箱:%s,验证码无效:%v", req.Account, req.Code)
		return
	}
	s.code.DelCode(req.Account)
	// 通过邮箱拉取数据库账号信息
	yostarAccount := db.GetYostarAccountByYostarAccount(req.Account)
	if yostarAccount == nil {
		if !config.GetAutoRegistration() {
			logger.Debug("邮箱:%s,账号不存在且关闭自动注册 account", req.Account)
			return
		}
		logger.Debug("邮箱:%s,账号不存在进行自动注册 account", req.Account)
		yostarAccount, err = db.AddYostarAccountByYostarAccount(req.Account)
		if err != nil {
			logger.Debug("自动注册sdk账号失败,数据库错误:%s account", err.Error())
			return
		}
	}
	if yostarAccount == nil {
		logger.Debug("邮箱:%s,进行数据库操作时候有未知错误 account", req.Account)
		return
	}
	// 更新token
	yostarAccount.YostarToken = fmt.Sprintf("%v-%s", yostarAccount.YostarUid, alg.RandStr(100))
	if err = db.UpdateYostarAccount(yostarAccount); err != nil {
		logger.Debug("数据库写入出现错误:%s account", err.Error())
		return
	}
	rsp.Result = 0
	rsp.YostarUid = strconv.Itoa(int(yostarAccount.YostarUid))
	rsp.YostarToken = yostarAccount.YostarToken
	rsp.YostarAccount = yostarAccount.YostarAccount
	logger.Debug("邮箱:%s,验证码登录成功 Code:%v,Token:%s,Uid:%v", req.Account, req.Code, yostarAccount.YostarToken, yostarAccount.YostarUid)
}
