package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCharacterPotentialStatExcel() {
	g.GetExcel().CharacterPotentialStatExcel = make([]*sro.CharacterPotentialStatExcel, 0)
	name := "CharacterPotentialStatExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CharacterPotentialStatExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().CharacterPotentialStatExcel))
}

type CharacterPotentialStat struct {
	CharacterPotentialList map[int64]map[int32]*sro.CharacterPotentialStatExcel
}

func (g *GameConfig) gppCharacterPotentialStatExcel() {
	g.GetGPP().CharacterPotentialStat = &CharacterPotentialStat{
		CharacterPotentialList: make(map[int64]map[int32]*sro.CharacterPotentialStatExcel),
	}
	for _, v := range g.GetExcel().GetCharacterPotentialStatExcel() {
		if g.GetGPP().CharacterPotentialStat.CharacterPotentialList[v.PotentialStatGroupId] == nil {
			g.GetGPP().CharacterPotentialStat.CharacterPotentialList[v.PotentialStatGroupId] = make(map[int32]*sro.CharacterPotentialStatExcel)
		}
		g.GetGPP().CharacterPotentialStat.CharacterPotentialList[v.PotentialStatGroupId][v.PotentialLevel] = v
	}

	logger.Info("处理角色能力解放升级配置完成,数量:%v个",
		len(g.GetGPP().CharacterPotentialStat.CharacterPotentialList))
}

func GetCharacterPotentialStatExcel(guid int64, level int32) *sro.CharacterPotentialStatExcel {
	if GC.GetGPP().CharacterPotentialStat.CharacterPotentialList[guid] == nil {
		return nil
	}
	return GC.GetGPP().CharacterPotentialStat.CharacterPotentialList[guid][level]
}
