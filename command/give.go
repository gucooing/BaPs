package command

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/protocol/proto"
)

type ApiGive struct {
	Uid      int64   `json:"uid"`
	ItemList []*Item `json:"item_list"`
}

type Item struct {
	Type int32 `json:"type"`
	Num  int64 `json:"num"`
	Id   int64 `json:"id"`
}

func (c *Command) Give(g *gin.Context) {
	req := new(ApiGive)
	err := g.BindJSON(req)
	if err != nil {
		g.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "ApiGive 解析错误",
		})
		return
	}
	s := enter.GetSessionByAccountServerId(req.Uid)
	if s == nil {
		g.JSON(http.StatusOK, gin.H{
			"code": 2,
			"msg":  "玩家不在线",
		})
		return
	}
	mail := &sro.MailInfo{
		Sender:         "gucooing",
		Comment:        "请查收您的意外奖励",
		SendDate:       time.Now().Unix(),
		ExpireDate:     time.Now().Add(24 * time.Hour).Unix(),
		ParcelInfoList: GiveJsonToProtobuf(req.ItemList),
	}
	if game.AddMail(s, mail) {
		g.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "发送成功",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"code": 3,
		"msg":  "邮件发送失败",
	})
}

func GiveJsonToProtobuf(bin []*Item) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, v := range bin {
		if _, ok := proto.ParcelType_name[v.Type]; !ok {
			continue
		}
		list = append(list, &sro.ParcelInfo{
			Type: v.Type,
			Id:   v.Id,
			Num:  v.Num,
		})
	}
	return list
}
