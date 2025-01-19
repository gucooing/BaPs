package command

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/pkg/alg"
)

type Command struct {
	router *gin.Engine
}

func NewCommand(router *gin.Engine) {
	c := &Command{
		router: router,
	}
	c.Router()
}

func (c *Command) Router() {
	if c == nil {
		return
	}
	gucooingApi := c.router.Group("/gucooing/api", alg.AutoGucooingApi())
	{
		gucooingApi.POST("/give", c.Give)
		gucooingApi.POST("/give_all", c.GiveAll)
		gucooingApi.GET("/getPlayerBin", c.getPlayerBin)
		gucooingApi.POST("/mail", c.Mail)
	}
}
