//go:build dev

package gdconf

import (
	"fmt"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
)

func (g *GameConfig) LoadExcel() {
	// 验证文件夹是否存在
	g.excelPath = g.resPath + "/Excel"
	dirInfo, err := os.Stat(g.excelPath)
	if err != nil || !dirInfo.IsDir() {
		info := fmt.Sprintf("找不到文件夹:%s,err:%s", g.excelPath, err)
		panic(info)
	}
	g.excelPath += "/"

	// 初始化excel
	g.Excel = new(sro.Excel)
	g.loadFunc = []func(){
		g.loadCafeInfoExcelTable,
		g.loadDefaultCharacterExcelTable,
		g.loadCharacterExcelTable,
	}

	for _, fn := range g.loadFunc {
		fn()
	}
}
