package sdk

import (
	"github.com/gin-gonic/gin"
)

// 通过邮箱拉取验证码
func (s *SDK) getEmailCode(c *gin.Context) {
	account := c.Query("account")
	code := s.code.GetCode(account)
	c.JSON(200, gin.H{
		"account": account,
		"code":    code,
	})
}
