package command

import (
	"errors"

	"github.com/gucooing/BaPs/common/code"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandGetEmailCode() {
	getEmailCode := &cdq.Command{
		Name:        "getemailcode",
		AliasList:   []string{"getemailcode", "gec"},
		Description: "获取指定邮箱的登录验证码",
		Permissions: cdq.Admin,
		Options: []*cdq.CommandOption{
			{
				Name:        "account",
				Description: "目标邮箱",
				Required:    true,
			},
		},
		CommandFunc: c.getEmailCode,
	}

	c.c.ApplicationCommand(getEmailCode)
}

// 通过邮箱拉取验证码
func (c *Command) getEmailCode(options map[string]*cdq.CommandOption) (string, error) {
	accountOption, ok := options["account"]
	if !ok {
		return "", errors.New("缺少参数 account")
	}
	if codeInfo := code.GetCodeInfo(accountOption.Option); codeInfo != nil &&
		codeInfo.FialNum < code.MaxFialNum {
		return string(codeInfo.Code), nil
	} else {
		return "", errors.New("验证码已过期或失效")
	}
}
