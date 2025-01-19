package sdk

import (
	"github.com/gin-gonic/gin"
	cd "github.com/gucooing/BaPs/common/code"
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
