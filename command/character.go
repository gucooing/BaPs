package command

import (
	"errors"
	"fmt"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandCharacter() {
	character := &cdq.Command{
		Name:        "character",
		AliasList:   []string{"character", "c"},
		Description: "直接设置角色,角色不存在时将添加角色",
		Permissions: cdq.User,
		Options: append(playerOptions, []*cdq.CommandOption{
			{
				Name:        "id",
				Description: "角色id",
				Required:    true,
			},
			{
				Name:        "level",
				Description: "要设置的角色类型",
				Required:    false,
			},
			{
				Name:        "starGrade",
				Description: "要设置的角色星级",
				Required:    false,
			},
			{
				Name:        "favorRank",
				Description: "要设置的角色好感度等级",
				Required:    false,
			},
			{
				Name:        "max",
				Description: "角色属性设置满级",
				Required:    false,
			},
		}...),
		CommandFunc: c.character,
	}

	c.c.ApplicationCommand(character)
}

func (c *Command) character(options map[string]*cdq.CommandOption) (string, error) {
	uidOption, ok := options["uid"]
	if !ok {
		return "", errors.New("缺少参数 uid")
	}
	idOption, ok := options["id"]
	if !ok {
		return "", errors.New("缺少参数 id")
	}

	// 玩家验证
	uid := alg.S2I64(uidOption.Option)
	s := enter.GetSessionByAccountServerId(uid)
	if s == nil {
		return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
	}
	respMsg := ""

	characterInfo := game.GetCharacterInfo(s, alg.S2I64(idOption.Option))
	if characterInfo == nil {
		if game.AddCharacter(s, alg.S2I64(idOption.Option)) {
			characterInfo = game.GetCharacterInfo(s, alg.S2I64(idOption.Option))
		} else {
			return "", errors.New("角色id错误")
		}
	}
	if _, ok := options["max"]; ok {
		if game.SetMaxCharacter(characterInfo) {
			return "已设置角色满级", nil
		}
	}
	if option, ok := options["level"]; ok {
		if game.SetCharacterLevel(characterInfo, alg.S2I32(option.Option)) {
			respMsg += "角色等级设置成功|"
		}
	}
	if option, ok := options["starGrade"]; ok {
		if game.SetCharacterStarGrade(characterInfo, alg.S2I32(option.Option)) {
			respMsg += "角色星级设置成功|"
		}
	}
	if option, ok := options["favorRank"]; ok {
		if game.SetCharacterFavorRank(characterInfo, alg.S2I32(option.Option)) {
			respMsg += "角色好感度等级设置成功|"
		}
	}

	return respMsg, nil
}
