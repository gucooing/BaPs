package gdconf

import (
	"encoding/json"
	"github.com/gucooing/BaPs/pkg/logger"
	"os"
)

type MailInfo struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

func (g *GameConfig) loadMailInfo() {
	g.GetGPP().MailInfoMap = make(map[string]*MailInfo)
	name := "Mail.json"
	file, err := os.ReadFile(g.dataPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetGPP().MailInfoMap); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("邮件配置读取成功文件:%s 读取成功,解析数量:%v", name, len(g.GetGPP().MailInfoMap))
}

func GetMailInfo(name string) *MailInfo {
	return GC.GPP.MailInfoMap[name]
}
