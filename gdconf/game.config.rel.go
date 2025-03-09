//go:build !dev
// +build !dev

package gdconf

import (
	"fmt"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	pb "google.golang.org/protobuf/proto"
)

func (g *GameConfig) LoadExcel() {
	dirInfo, err := os.Stat(g.dataPath)
	if err != nil || !dirInfo.IsDir() {
		info := fmt.Sprintf("找不到文件夹:%s,err:%s", g.dataPath, err)
		panic(info)
	}
	g.dataPath += "/"

	file, err := os.ReadFile(g.dataPath + "Excel.bin")
	if err != nil {
		logger.Error("Excel.bin 读取失败,err:%s", err)
		return
	}
	bin, err := mx.DeExcelBytes(file, 1618496251562615)
	if err != nil {
		logger.Error("解析Excel失败")
		return
	}
	g.Excel = new(sro.Excel)
	err = pb.Unmarshal(bin, g.Excel)
	if err != nil {
		logger.Error("解析Excel失败,err:%s", err)
		return
	}
}
