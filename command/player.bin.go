package command

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandGetPlayer() {
	getPlayer := &cdq.Command{
		Name:        "getplayer",
		AliasList:   []string{"getplayer", "gp"},
		Description: "获取指定玩家的db数据",
		Permissions: cdq.Admin,
		Options: append(playerOptions, &cdq.CommandOption{
			Name:        "bin",
			Description: "是否输出二进制pb数据,默认:0",
			Required:    false,
		}),
		CommandFunc: c.getPlayer,
	}

	c.c.ApplicationCommand(getPlayer)
}

func (c *Command) getPlayer(options map[string]*cdq.CommandOption) (string, error) {
	uidOption, ok := options["uid"]
	if !ok {
		return "", errors.New("缺少参数 uid")
	}
	isBin := int64(0)
	isBinOption, ok := options["bin"]
	if ok {
		isBin = alg.MaxInt64(alg.S2I64(isBinOption.Option), 0)
	}

	uid := alg.S2I64(uidOption.Option)
	session := enter.GetSessionByUid(uid)
	if session == nil {
		return "", errors.New(fmt.Sprintf("玩家未注册 UID:%v", uid))
	}

	if isBin == 0 {
		return sonic.MarshalString(session.PlayerBin)
	} else {
		return hex.EncodeToString(session.GetPbBinData()), nil
	}
}
