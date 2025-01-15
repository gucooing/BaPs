package sdk

import (
	"github.com/gin-gonic/gin"
	cd "github.com/gucooing/BaPs/common/code"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/alg"
)

// 通过邮箱拉取验证码
func (s *SDK) getEmailCode(c *gin.Context) {
	account := c.Query("account")
	var code int32 = 0
	var msg string
	defer c.JSON(200, gin.H{
		"account": account,
		"code":    code,
		"msg":     msg,
	})
	if codeInfo := s.code.GetCodeInfo(account); codeInfo != nil &&
		codeInfo.FialNum < cd.MaxFialNum {
		code = codeInfo.Code
	} else {
		msg = "验证码已过期或失效"
	}
}

func (s *SDK) getPlayerBin(c *gin.Context) {
	uid := c.Query("uid")

	session := enter.GetSessionByAccountServerId(alg.S2I64(uid))
	if session == nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "player err",
		})
		return
	}
	c.JSON(200, session)
}
