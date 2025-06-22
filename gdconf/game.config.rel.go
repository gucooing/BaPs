//go:build !debug
// +build !debug

package gdconf

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
	pb "google.golang.org/protobuf/proto"
)

var excelPath string
var excelDbPath string

//go:embed Excel.bin
var excelBytes []byte

func (g *GameConfig) LoadExcel() {
	//ty:
	dirInfo, err := os.Stat(g.dataPath)
	if err != nil || !dirInfo.IsDir() {
		info := fmt.Sprintf("找不到文件夹:%s,err:%s", g.dataPath, err)
		panic(info)
	}
	g.dataPath += "/"

	//excelBytes, err := os.ReadFile(g.dataPath + "Excel.bin")
	//if err != nil {
	//	if os.IsNotExist(err) {
	//		logger.Error("没有找到Excel.bin尝试自动下载....")
	//		err := g.downloadExcel(g.dataPath + "Excel.bin")
	//		if err == nil {
	//			logger.Error("Excel.bin自动下载成功！")
	//			goto ty
	//		}
	//	}
	//	logger.Error("Excel.bin 读取失败,err:%s", err)
	//	return
	//}
	bin, err := mx.DeExcelBytes(excelBytes)
	if err != nil {
		panic("Excel.bin不匹配")
		return
	}
	g.Excel = new(sro.Excel)
	err = pb.Unmarshal(bin, g.Excel)
	if err != nil {
		panic("解析Excel失败,请检查Excel.bin版本和服务端版本是否一致")
		return
	}
}

func (g *GameConfig) downloadExcel(path string) error {
	resp, err := http.Get(g.excelUrl)
	if err != nil {
		logger.Error("下载Excel.bin失败,请手动下载")
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func loadExcelFile[T any](path string, table *[]*T) {
	*table = make([]*T, 0)
}
