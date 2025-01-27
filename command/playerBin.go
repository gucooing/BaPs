package command

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/code"
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/alg"
	pb "google.golang.org/protobuf/proto"
)

func (c *Command) getPlayerBin(g *gin.Context) {
	uid := g.Query("uid")

	if session := enter.GetSessionByAccountServerId(alg.S2I64(uid)); session != nil {
		session.GoroutinesSync.Lock()
		defer session.GoroutinesSync.Unlock()
		g.JSON(200, session)
		return
	}
	yostarGame := db.GetYostarGameByAccountServerId(alg.S2I64(uid))
	if yostarGame != nil {
		playerBin := new(sro.PlayerBin)
		if pb.Unmarshal(yostarGame.BinData, playerBin) == nil {
			g.JSON(200, playerBin)
			return
		}
	}
	g.JSON(200, gin.H{
		"code": 2,
		"msg":  "player err",
	})
}

// 通过邮箱拉取验证码
func (c *Command) getEmailCode(g *gin.Context) {
	account := g.Query("account")
	var rspCode int32 = 0
	var msg string
	defer g.JSON(200, gin.H{
		"account": account,
		"code":    rspCode,
		"msg":     msg,
	})
	if codeInfo := code.GetCodeInfo(account); codeInfo != nil &&
		codeInfo.FialNum < code.MaxFialNum {
		rspCode = codeInfo.Code
	} else {
		msg = "验证码已过期或失效"
	}
}
