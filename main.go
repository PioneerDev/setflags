package main

import (
	"fmt"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/setting"
	"set-flags/routers"
)

func main() {

	models.InitDB()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        routers.InitRouter(),
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
