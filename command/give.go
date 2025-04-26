package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandGive() {
	give := &cdq.Command{
		Name:        "give",
		AliasList:   []string{"give", "g"},
		Description: "获取物品,注:一次只能获取一个物品",
		Permissions: cdq.User,
		Options: append(playerOptions, []*cdq.CommandOption{
			{
				Name:        "id",
				Description: "需要给予物品的id",
				Required:    true,
			},
			{
				Name:        "t",
				Description: "需要给予物品的类型",
				Required:    true,
			},
			{
				Name:        "num",
				Description: "需要给予物品的数量 默认值:1",
				Required:    false,
			},
		}...),
		CommandFunc: c.give,
	}

	c.c.ApplicationCommand(give)
}

func (c *Command) give(options map[string]*cdq.CommandOption) (string, error) {
	uidOption, ok := options["uid"]
	if !ok {
		return "", errors.New("缺少参数 uid")
	}
	idOption, ok := options["id"]
	if !ok {
		return "", errors.New("缺少参数 id")
	}
	typeOption, ok := options["t"]
	if !ok {
		return "", errors.New("缺少参数 t")
	}
	parcelType := proto.GetParcelTypeValue(typeOption.Option)
	num := int64(1)
	itemNum, ok := options["num"]
	if ok {
		num = alg.MaxInt64(alg.S2I64(itemNum.Option), 1)
	}

	// 玩家验证
	uid := alg.S2I64(uidOption.Option)
	s := enter.GetSessionByAccountServerId(uid)
	if s == nil {
		return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
	}

	// 执行
	s.GoroutinesSync.Lock()
	defer s.GoroutinesSync.Unlock()

	mail := &sro.MailInfo{
		Sender:     "gucooing",
		Comment:    "请查收您的意外奖励",
		SendDate:   time.Now().Unix(),
		ExpireDate: time.Now().Add(10 * time.Minute).Unix(),
		ParcelInfoList: []*sro.ParcelInfo{
			{
				Type: parcelType.Value(),
				Id:   alg.S2I64(idOption.Option),
				Num:  num,
			},
		},
	}
	if game.AddMail(s, mail) {
		return "请查询游戏内邮箱获取结果", nil
	}
	return "", errors.New("执行give 失败,游戏邮箱错误")
}
