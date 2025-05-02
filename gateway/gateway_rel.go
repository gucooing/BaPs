//go:build !dev
// +build !dev

package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/protocol/mx"
)

func status(router *gin.Engine) {

}

func logPlayerMsg(logType int, msg mx.Message) {

}
