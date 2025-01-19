package command

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
)

type ApiMail struct {
	Player         bool   `json:"player"`           // 是否私人邮件
	Uid            int64  `json:"uid"`              // 玩家id
	Sender         string `json:"sender"`           // 发人
	Comment        string `json:"comment"`          // 内容
	SendDate       int64  `json:"send_date"`        // 发时
	ExpireDate     int64  `json:"expire_date"`      // 截至时
	ParcelInfoList string `json:"parcel_info_list"` // 附件 json
	DelMail        bool   `json:"del_mail"`         // 是否删除邮件
	MailId         int64  `json:"mail_id"`          // 邮件id
}

func (c *Command) Mail(g *gin.Context) {
	req := new(ApiMail)
	if err := g.ShouldBind(req); err != nil {
		g.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "ApiMail 解析错误",
		})
		return
	}
	if req.DelMail {
		c.DelMail(g, req)
		return
	}
	if req.Player {
	}
}

func (c *Command) DelMail(g *gin.Context, info *ApiMail) {
	if info.Player {
		s := enter.GetSessionByAccountServerId(info.Uid)
		if s == nil {
			g.JSON(http.StatusOK, gin.H{
				"code": 2,
				"msg":  "玩家不在线",
			})
			return
		}
		if game.DelMail(s, info.MailId) {
			g.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "删除邮件成功",
			})
			return
		}
		g.JSON(http.StatusOK, gin.H{
			"code": 3,
			"msg":  "删除邮件失败,邮件不存在或者玩家不在线",
		})
		return
	}
}
