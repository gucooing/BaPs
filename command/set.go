package command

import (
	"errors"
	"fmt"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/cdq"
)

type ApiSet struct {
	s    *enter.Session
	Type string
	Sub1 string
}

func (c *Command) ApplicationCommandSet() {
	set := &cdq.Command{
		Name:        "set",
		AliasList:   []string{"set"},
		Description: "直接设置一个值",
		Permissions: cdq.User,
		Options: append(playerOptions, []*cdq.CommandOption{
			{
				Name:        "type",
				Description: "设置类型",
				Required:    true,
			},
			{
				Name:        "sub1",
				Description: "设置参数1",
				Required:    true,
			},
		}...),
		CommandFunc: c.set,
	}

	c.c.ApplicationCommand(set)
}

func (c *Command) set(options map[string]*cdq.CommandOption) (string, error) {
	uidOption, ok := options["uid"]
	if !ok {
		return "", errors.New("缺少参数 uid")
	}
	typeOption, ok := options["type"]
	if !ok {
		return "", errors.New("缺少参数 type")
	}
	sub1Option, ok := options["sub1"]
	if !ok {
		return "", errors.New("缺少参数 sub1")
	}

	// 玩家验证
	uid := alg.S2I64(uidOption.Option)
	s := enter.GetSessionByAccountServerId(uid)
	if s == nil {
		return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
	}
	s.GoroutinesSync.Lock()
	defer s.GoroutinesSync.Unlock()

	req := &ApiSet{
		s:    s,
		Type: typeOption.Option,
		Sub1: sub1Option.Option,
	}

	return req.parseSet()
}

func (x *ApiSet) parseSet() (string, error) {
	switch x.Type {
	case "AccountLevel":
		return x.setAccountLevel()
	case "Toast":
		return x.setToast()
	default:
		return "", errors.New(fmt.Sprintf("Set Type 未实现:%s", x.Type))
	}
}

func (x *ApiSet) setAccountLevel() (string, error) {
	game.SetAccountLevel(x.s, alg.S2I32(x.Sub1))
	game.AddToast(x.s, "已设置账号等级,请重新登录以刷新")
	return fmt.Sprintf("已设置玩家等级:%v级", game.GetAccountLevel(x.s)), nil
}

func (x *ApiSet) setToast() (string, error) {
	game.AddToast(x.s, x.Sub1)
	return fmt.Sprintf("已设置玩家通知:%s", x.Sub1), nil
}
