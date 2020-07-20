package main

import (
	"fmt"
	"log"
	"set-flags/models"
	"set-flags/pkg/logging"
	"set-flags/pkg/setting"
	"set-flags/routers"
	"syscall"
	"time"

	"github.com/fvbock/endless"
)

func main() {
	logging.Setup()
	models.InitDB()

	endless.DefaultReadTimeOut = time.Duration(setting.GetConfig().App.ReadTimeOut) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(setting.GetConfig().App.WriteTimeOut) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.GetConfig().App.HTTPPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		logging.Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		logging.Error(fmt.Sprintf("Server err: %v", err))
		log.Printf("Server err: %v", err)
	}
}
