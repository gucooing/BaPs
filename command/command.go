package command

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/cdq"
	cdqlog "github.com/gucooing/cdq/logger"
)

type Command struct {
	c *cdq.CDQ
}

func NewCommand(router *gin.Engine) {
	command := new(Command)
	command.c = cdq.New(&cdq.CDQ{Log: cdqlog.NewLog(cdqlog.LevelInfo, nil)})
	ginApi := cdq.NewGinApi(command.c)
	ginApi.SetRouter(router)
	ginApi.SetApiKey(config.GetGucooingApiKey())
	command.c.AddCommandRun(ginApi)

	// 注册指令
	command.ApplicationCommandGiveAll()
	command.ApplicationCommandGive()
	command.ApplicationCommandGetPlayer()
	command.ApplicationCommandGetEmailCode()
	command.ApplicationCommandMail()
	command.ApplicationCommandSet()
	command.ApplicationCommandPing()
	command.ApplicationCommandCharacter()
}

var playerOptions = []*cdq.CommandOption{
	{
		Name:        "uid",
		Description: "玩家游戏id",
		Required:    true,
	},
}
