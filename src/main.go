package main

import (
	"fmt"
	"log"
	"set-flags/src/global"
	"set-flags/src/models"
	"set-flags/src/pkg/setting"
	"set-flags/src/routers"
	"syscall"
	"time"

	"github.com/fvbock/endless"
)

func main() {
	models.InitDB()
	global.BotInit()

	endless.DefaultReadTimeOut = time.Duration(setting.GetConfig().App.ReadTimeOut) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(setting.GetConfig().App.WriteTimeOut) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.GetConfig().App.HTTPPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d\n", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v\n", err)
	}
}
