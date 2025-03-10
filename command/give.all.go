package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
	"github.com/gucooing/cdq"
)

type ApiGiveAll struct {
	Uid  int64  `json:"uid"`
	Type string `json:"type"`
	Num  int64  `json:"num"`
}

func (c *Command) ApplicationCommandGiveAll() {
	giveAll := &cdq.Command{
		Name:        "giveall",
		AliasList:   []string{"giveall", "ga"},
		Description: "获取某个类型的全部物品",
		Permissions: cdq.User,
		Options: append(playerOptions, []*cdq.CommandOption{
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
		CommandFunc: c.giveALL,
	}

	c.c.ApplicationCommand(giveAll)
}

func (c *Command) giveALL(options map[string]*cdq.CommandOption) (string, error) {
	uidOption, ok := options["uid"]
	if !ok {
		return "", errors.New("缺少参数 uid")
	}
	typeOption, ok := options["t"]
	if !ok {
		return "", errors.New("缺少参数 t")
	}
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
	parcelInfoList := GiveAllJsonToProtobuf(&ApiGiveAll{
		Uid:  uid,
		Type: typeOption.Option,
		Num:  num,
	})
	if len(parcelInfoList) == 0 {
		return "", errors.New(fmt.Sprintf("不存在此物品类型 Type:%s", typeOption.Option))
	}

	mail := &sro.MailInfo{
		Sender:         "gucooing",
		Comment:        "请查收您的意外奖励",
		SendDate:       time.Now().Unix(),
		ExpireDate:     time.Now().Add(10 * time.Minute).Unix(),
		ParcelInfoList: parcelInfoList,
	}
	if game.AddMail(s, mail) {
		return "请查询游戏内邮箱获取结果", nil
	}
	return "", errors.New("游戏邮箱错误")
}

func GiveAllJsonToProtobuf(req *ApiGiveAll) []*sro.ParcelInfo {
	switch req.Type {
	case "Material": // 材料
		return GiveAllMaterial(req)
	case "Character": // 角色
		return GiveAllCharacter(req)
	case "Equipment": // 装备
		return GiveAllEquipment(req)
	case "Furniture": // 家具
		return GiveAllFurniture(req)
	case "Favor": // 礼物
		return GiveAllFavor(req)
	case "Emblem": // 称号
		return GiveAllEmblem(req)
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

func GiveAllFurniture(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetFurnitureExcelTableMap() {
		num := req.Num
		if num <= 0 {
			num = 2
		}
		list = append(list, &sro.ParcelInfo{
			Type: proto.ParcelType_Furniture,
			Id:   conf.Id,
			Num:  num,
		})
	}

	return list
}

func GiveAllFavor(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	if req.Num <= 0 {
		req.Num = 999
	}
	for _, conf := range gdconf.GetItemExcelCategoryMap("Favor") {
		list = append(list, &sro.ParcelInfo{
			Type: proto.ParcelType_Item,
			Id:   conf.Id,
			Num:  req.Num,
		})
	}

	return list
}

func GiveAllEmblem(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetEmblemExcelList() {
		list = append(list, &sro.ParcelInfo{
			Type: proto.ParcelType_Emblem,
			Id:   conf.Id,
		})
	}
	return list
}
