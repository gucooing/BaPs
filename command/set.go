package command

import (
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
		s := enter.GetSessionByAccountServerId(req.Uid)
		if s == nil {
			g.JSON(http.StatusOK, gin.H{
				"code": 2,
				"msg":  "玩家不在线",
			})
			return
		}
		game.SetAccountLevel(s, alg.S2I32(req.Sub1))
	default:
		g.JSON(http.StatusOK, gin.H{
			"code": 4,
			"msg":  "Set Type 失败",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}
