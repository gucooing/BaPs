package command

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/cdq"
)

type ApiMail struct {
	Player         bool   `json:"player"`           // 是否私人邮件
	Uid            int64  `json:"uid"`              // 玩家id
	Sender         string `json:"sender"`           // 发件人
	Comment        string `json:"comment"`          // 内容
	SendDate       int64  `json:"send_date"`        // 发时
	ExpireDate     int64  `json:"expire_date"`      // 截至时
	ParcelInfoList string `json:"parcel_info_list"` // 附件 json
}

func (c *Command) ApplicationCommandMail() {
	mail := &cdq.Command{
		Name:        "mail",
		AliasList:   []string{"mail", "m"},
		Description: "发送一封邮件,对象可以是全局或单个玩家",
		Permissions: cdq.Admin,
		Options: []*cdq.CommandOption{
			{
				Name:        "player",
				Description: "是否为玩家邮件,0为全局邮件",
				Required:    true,
			},
			{
				Name:        "sender",
				Description: "发件人",
				Required:    true,
			},
			{
				Name:        "comment",
				Description: "邮件文本内容",
				Required:    true,
			},
			{
				Name:        "send_date",
				Description: "领取开始时间",
				Required:    true,
			},
			{
				Name:        "expire_date",
				Description: "邮件有效截至时间",
				Required:    true,
			},
			{
				Name:        "uid",
				Description: "如果是玩家邮件,则此选项必写",
				Required:    false,
			},
			{
				Name:        "parcel_info_list",
				Description: "附件 json格式",
				Required:    false,
			},
		},
		CommandFunc: c.mail,
	}

	c.c.ApplicationCommand(mail)
}

func (c *Command) mail(options map[string]*cdq.CommandOption) (string, error) {
	playerOption, ok := options["player"]
	if !ok {
		return "", errors.New("缺少参数 player")
	}
	senderOption, ok := options["sender"]
	if !ok {
		return "", errors.New("缺少参数 sender")
	}
	commentOption, ok := options["comment"]
	if !ok {
		return "", errors.New("缺少参数 comment")
	}
	sendDateOption, ok := options["send_date"]
	if !ok {
		return "", errors.New("缺少参数 send_date")
	}
	expireDateOption, ok := options["expire_date"]
	if !ok {
		return "", errors.New("缺少参数 expire_date")
	}

	str := ""
	if parcelInfoListOption, ok := options["parcel_info_list"]; ok {
		str = parcelInfoListOption.Option
	}

	if alg.S2I32(playerOption.Option) == 0 {
		parcelInfoList, err := genParcelInfo[db.ParcelInfo](str)
		if err != nil {
			return "", errors.New(fmt.Sprintf("解析邮件附件失败:%s", err.Error()))
		}
		sendMail := &db.YostarMail{
			Sender:         senderOption.Option,
			Comment:        commentOption.Option,
			SendDate:       sql.NullTime{Time: time.Unix(alg.S2I64(sendDateOption.Option), 0)},
			ExpireDate:     sql.NullTime{Time: time.Unix(alg.S2I64(expireDateOption.Option), 0)},
			ParcelInfoList: parcelInfoList,
		}
		if enter.AddYostarMail(sendMail) {
			return "全局邮件发送成功", nil
		}
		return "", errors.New("全局邮件发送失败")
	} else {
		parcelInfoList, err := genParcelInfo[sro.ParcelInfo](str)
		if err != nil {
			return "", errors.New(fmt.Sprintf("解析邮件附件失败:%s", err.Error()))
		}
		sendMail := &sro.MailInfo{
			Sender:         senderOption.Option,
			Comment:        commentOption.Option,
			SendDate:       alg.S2I64(sendDateOption.Option),
			ExpireDate:     alg.S2I64(expireDateOption.Option),
			ParcelInfoList: parcelInfoList,
		}
		uidOption, ok := options["uid"]
		if !ok {
			return "", errors.New("缺少参数 uid")
		}
		// 玩家验证
		uid := alg.S2I64(uidOption.Option)
		s := enter.GetSessionByAccountServerId(uid)
		if s == nil {
			return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
		}
		s.GoroutinesSync.Lock()
		defer s.GoroutinesSync.Unlock()
		if game.AddMail(s, sendMail) {
			return "请查询游戏内邮箱获取结果", nil
		}
		return "", errors.New("私人邮件发送失败")
	}
}

func genParcelInfo[T any](str string) ([]*T, error) {
	list := make([]*T, 0)
	if str == "" {
		return list, nil
	}
	err := json.Unmarshal([]byte(str), &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
