package main

import (
	"log"

	"github.com/gucooing/BaPs"
)

func main() {
	log.Print("BaPs动态库加载成功")
}

//export StartServer
func StartServer() {
	BaPs.NewBaPs()
}
