//go:build dev
// +build dev

package gateway

import (
	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
	"github.com/bytedance/sonic"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/cmd"
	"github.com/gucooing/BaPs/protocol/mx"
)

func status(router *gin.Engine) {
	go example.Work()
	pprof.Register(router, "pprof")
	srv, _ := statsviz.NewServer()
	router.GET("/debug/statsviz/*filepath", func(context *gin.Context) {
		if context.Param("filepath") == "/ws" {
			srv.Ws()(context.Writer, context.Request)
			return
		}
		srv.Index()(context.Writer, context.Request)
	})
}

func logPlayerMsg(logType int, msg mx.Message) {
	cmdId := cmd.Get().GetCmdIdByProtoObj(msg)
	if _, ok := config.GetBlackCmd()[cmdId.String()]; ok ||
		!config.GetIsLogMsgPlayer() {
		return
	}
	var a string
	switch logType {
	case Client:
		a = "@LogTag(player_msg)@ gateway c--->s cmd id:"
	case Server:
		a = "@LogTag(player_msg)@ gateway s--->c cmd id:"
	case NoRoute:
		a = "@LogTag(player_no_route)@ c --> s no route for msg, cmd id:"
	}
	b, _ := sonic.MarshalIndent(msg, "", "  ")

	logger.Debug("%s%s :%s", a, cmdId.String(), string(b))
}
