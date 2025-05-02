package mail

import (
	"errors"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/wneessen/go-mail"
)

var mc *Mail

type Mail struct {
	*mail.Client
	from string
}

func NewMail() {
	conf := config.GetMail()
	if conf == nil {
		return
	}
	logger.Info("检测到邮件服务配置,正在尝试连接")

	client, err := mail.NewClient(conf.Host, mail.WithTLSPortPolicy(mail.TLSOpportunistic), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(conf.Username), mail.WithPassword(conf.Password))
	if err != nil {
		logger.Error("连接邮件服务失败")
		return
	}
	logger.Info("邮件服务连接成功")

	mc = &Mail{
		Client: client,
		from:   conf.Username,
	}
}

func SendTextMail(username, ject, body string) error {
	if mc == nil {
		return nil
	}
	message := mail.NewMsg()
	if err := message.From(mc.from); err != nil {
		return errors.New("邮件服务配置错误")
	}
	if err := message.To(username); err != nil {
		return errors.New("收件人不合法")
	}
	message.Subject(ject)
	message.SetBodyString(mail.TypeTextPlain, body)

	if err := mc.DialAndSend(message); err != nil {
		return errors.New("邮件发送失败")
	}
	return nil
}
