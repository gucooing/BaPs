package gdconf

import (
	"encoding/json"
	"math/rand"
	"os"

	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

type ArenaNPCInfo struct {
	ArenaNPCList []*ArenaNPC
	ArenaNPCMap  map[int64]*ArenaNPC
}

type ArenaNPC struct {
	Index                        int64   `json:"Index"`
	NpcaccountLevel              int64   `json:"NpcaccountLevel"`
	Npclevel                     int32   `json:"Npclevel"`
	NpclevelDeviation            int32   `json:"NpclevelDeviation"`
	NpcstarGrade                 int32   `json:"NpcstarGrade"`
	ExceptionMainCharacterIds    []int64 `json:"ExceptionMainCharacterIds"`
	ExceptionSupportCharacterIds []int64 `json:"ExceptionSupportCharacterIds"`
}

func (g *GameConfig) loadArenaNPC() {
	g.GetGPP().ArenaNPCInfo = &ArenaNPCInfo{
		ArenaNPCList: make([]*ArenaNPC, 0),
		ArenaNPCMap:  make(map[int64]*ArenaNPC),
	}
	name := "ArenaNPC.json"
	file, err := os.ReadFile(g.dataPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetGPP().ArenaNPCInfo.ArenaNPCList); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	for _, v := range g.GetGPP().ArenaNPCInfo.ArenaNPCList {
		g.GetGPP().ArenaNPCInfo.ArenaNPCMap[v.Index] = v
	}

	logger.Info("竞技场人机读取成功文件:%s 读取成功,解析数量:%v", name, len(g.GetGPP().ArenaNPCInfo.ArenaNPCList))
}

var DefaultArenaNPCInfo = &ArenaNPC{
	Index:                        0,
	NpcaccountLevel:              20,
	Npclevel:                     20,
	NpclevelDeviation:            3,
	NpcstarGrade:                 2,
	ExceptionMainCharacterIds:    []int64{10079, 10080},
	ExceptionSupportCharacterIds: []int64{20007, 26011},
}

// RandGetArenaNPC 随机取一个
func RandGetArenaNPC() *ArenaNPC {
	list := GC.GetGPP().ArenaNPCInfo.ArenaNPCList
	if len(list) == 0 {
		return DefaultArenaNPCInfo
	}
	return list[rand.Intn(len(list))]
}

// GetArenaNPCByIndex 通过index取得npc信息
func GetArenaNPCByIndex(index int64) *ArenaNPC {
	confs := GC.GetGPP().ArenaNPCInfo.ArenaNPCMap
	if conf, ok := confs[index]; ok {
		return conf
	}
	return DefaultArenaNPCInfo
}

func (x *ArenaNPC) GetArenaCharacterDB(id int64) *proto.ArenaCharacterDB {
	info := &proto.ArenaCharacterDB{
		UniqueId:               id,
		StarGrade:              x.NpcstarGrade,
		Level:                  x.Npclevel + rand.Int31n(x.NpclevelDeviation),
		PublicSkillLevel:       1,
		ExSkillLevel:           1,
		PassiveSkillLevel:      1,
		ExtraPassiveSkillLevel: 1,
		LeaderSkillLevel:       1,
	}

	return info
}
