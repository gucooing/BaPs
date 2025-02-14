package command

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
)

type ApiSet struct {
	Uid  int64  `json:"uid"`
	Type string `json:"type"`
	Sub1 string `json:"sub1"`
	Su   bool   `json:"-"`
	Msg  string `json:"-"`
}

func (c *Command) Set(g *gin.Context) {
	req := new(ApiSet)
	err := g.BindJSON(req)
	if err != nil {
		g.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "ApiSet 解析错误",
		})
		return
	}
	parseSet(req, g)
}

func parseSet(req *ApiSet, g *gin.Context) {
	if req == nil {
		return
	}
	switch req.Type {
	case "AccountLevel":
		SetAccountLevel(req)
	case "Toast":
		SetToast(req)
	case "MaxCharacter":
		SetMaxCharacter(req)
	default:
		req.Msg = "Set Type 未实现"
	}
	if req.Su {
		g.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  req.Msg,
		})
	} else {
		g.JSON(http.StatusOK, gin.H{
			"code": 2,
			"msg":  req.Msg,
		})
	}
}

func SetAccountLevel(req *ApiSet) {
	s := enter.GetSessionByAccountServerId(req.Uid)
	if s == nil {
		req.Msg = "玩家不在线"
		return
	}
	game.SetAccountLevel(s, alg.S2I32(req.Sub1))
	game.AddToast(s, "已设置账号等级,请重新登录以刷新")
	req.Su = true
	req.Msg = fmt.Sprintf("已设置玩家等级:%v级", game.GetAccountLevel(s))
}

func SetToast(req *ApiSet) {
	s := enter.GetSessionByAccountServerId(req.Uid)
	if s == nil {
		req.Msg = "玩家不在线"
		return
	}
	game.AddToast(s, req.Sub1)
	req.Su = true
	req.Msg = fmt.Sprintf("已设置玩家通知:%s", req.Sub1)
}

func SetMaxCharacter(req *ApiSet) {
	s := enter.GetSessionByAccountServerId(req.Uid)
	if s == nil {
		req.Msg = "玩家不在线"
		return
	}
	game.MaxAllCharacter(s)
	game.AddToast(s, "角色已全部满级,请重新登录刷新")
	req.Su = true
	req.Msg = fmt.Sprintf("已设置角色满级")
}
