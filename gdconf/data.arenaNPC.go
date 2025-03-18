package gdconf

import (
	"encoding/json"
	"math/rand"
	"os"

	"github.com/gucooing/BaPs/pkg/logger"
)

type ArenaNPCInfo struct {
	NpcaccountLevel              int64   `json:"NpcaccountLevel"`
	Npclevel                     int64   `json:"Npclevel"`
	NpclevelDeviation            int64   `json:"NpclevelDeviation"`
	NpcstarGrade                 int64   `json:"NpcstarGrade"`
	ExceptionMainCharacterIds    []int64 `json:"ExceptionMainCharacterIds"`
	ExceptionSupportCharacterIds []int64 `json:"ExceptionSupportCharacterIds"`
}

func (g *GameConfig) loadArenaNPC() {
	g.GetGPP().ArenaNPCList = make([]*ArenaNPCInfo, 0)
	name := "ArenaNPC.json"
	file, err := os.ReadFile(g.dataPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetGPP().ArenaNPCList); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("竞技场人机读取成功文件:%s 读取成功,解析数量:%v", name, len(g.GetGPP().ArenaNPCList))
}

var DefaultArenaNPCInfo = &ArenaNPCInfo{
	NpcaccountLevel:              20,
	Npclevel:                     20,
	NpclevelDeviation:            3,
	NpcstarGrade:                 2,
	ExceptionMainCharacterIds:    []int64{10079, 10080},
	ExceptionSupportCharacterIds: []int64{20007, 26011},
}

// GetArenaNPCInfo 随机取一个
func GetArenaNPCInfo() *ArenaNPCInfo {
	list := GC.GetGPP().ArenaNPCList
	if len(list) == 0 {
		return nil
	}
	return list[rand.Intn(len(list))]
}
