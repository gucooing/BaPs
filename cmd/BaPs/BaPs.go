package main

import (
	"bufio"
	"log"
	"os"

	"github.com/gucooing/BaPs"
	"github.com/gucooing/BaPs/pkg/logger"
)

func main() {
	exit := func() {
		log.Printf("\n执行结束请输入任何键退出程序....")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			return
		}
	}
	defer func() {
		logger.CloseLogger() // 等logger打印完后再退出
		if err := recover(); err != nil {
			log.Println("\n程序异常退出,原因:")
			log.Println(err)
			exit()
		}
	}()
	BaPs.NewBaPs()
}
