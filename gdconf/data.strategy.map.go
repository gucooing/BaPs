package gdconf

import (
	"encoding/json"
	"os"

	"github.com/gucooing/BaPs/pkg/logger"
)

type StrategyMap struct {
	StrategyMap       string       `json:"StrategyMap"`
	LastEnemyEntityId int64        `json:"LastEnemyEntityId"`
	EnemyInfos        []*EnemyInfo `json:"EnemyInfos"`
	StrategyObjects   []*EnemyInfo `json:"StrategyObjects"`
}
type EnemyInfo struct {
	EntityId int64    `json:"EntityId"`
	Id       int64    `json:"Id"`
	Rotate   *Vector3 `json:"Rotate"`
	Location *Vector3 `json:"Location"`
}
type Vector3 struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
	Z float32 `json:"Z"`
}

func (g *GameConfig) loadStrategyMap() {
	g.GetGPP().StrategyMap = make(map[string]*StrategyMap)
	name := "StrategyMap.json"
	file, err := os.ReadFile(g.dataPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetGPP().StrategyMap); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("格子地图读取成功文件:%s 读取成功,解析数量:%v", name, len(g.GetGPP().StrategyMap))
}

func GetStrategyMap(str string) *StrategyMap {
	return GC.GetGPP().StrategyMap[str]
}
