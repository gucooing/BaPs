package command

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
)

type ApiGiveAll struct {
	Uid  int64  `json:"uid"`
	Type string `json:"type"`
	Num  int64  `json:"num"`
}

func (c *Command) GiveAll(g *gin.Context) {
	req := new(ApiGiveAll)
	err := g.BindJSON(req)
	if err != nil {
		g.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "ApiGiveAll 解析错误",
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
	parcelInfoList := GiveAllJsonToProtobuf(req)
	if len(parcelInfoList) == 0 {
		g.JSON(http.StatusOK, gin.H{
			"code": 3,
			"msg":  "不存在此物品类型",
		})
		return
	}
	mail := &sro.MailInfo{
		Sender:         "gucooing",
		Comment:        "请查收您的意外奖励",
		SendDate:       time.Now().Unix(),
		ExpireDate:     time.Now().Add(10 * time.Minute).Unix(),
		ParcelInfoList: parcelInfoList,
	}
	if game.AddMail(s, mail) {
		g.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "发送成功",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"code": 4,
		"msg":  "邮件发送失败",
	})
}

func GiveAllJsonToProtobuf(req *ApiGiveAll) []*sro.ParcelInfo {
	switch req.Type {
	case "Material": // 材料
		return GiveAllMaterial(req)
	case "Character": // 角色
		return GiveAllCharacter(req)
	case "Equipment": // 装备
		return GiveAllEquipment(req)
	}
	return nil
}

func GiveAllMaterial(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	if req.Num <= 0 {
		req.Num = 999
	}
	for _, conf := range gdconf.GetItemExcelCategoryMap("Material") {
		list = append(list, &sro.ParcelInfo{
			Type: proto.ParcelType_Item,
			Id:   conf.Id,
			Num:  req.Num,
		})
	}

	return list
}

func GiveAllCharacter(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetCharacterMap() {
		list = append(list, &sro.ParcelInfo{
			Type: proto.ParcelType_Character,
			Id:   conf.Id,
			Num:  1,
		})
	}
	return list
}

func GiveAllEquipment(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetEquipmentExcelTableMap() {
		num := req.Num
		if conf.MaxLevel < 10 {
			if num <= 0 {
				num = alg.MaxInt64(conf.StackableMax/10, 1)
			}
		} else {
			if num <= 0 {
				num = 20
			}
		}

		list = append(list, &sro.ParcelInfo{
			Type: proto.ParcelType_Equipment,
			Id:   conf.Id,
			Num:  num,
		})
	}

	return list
}
