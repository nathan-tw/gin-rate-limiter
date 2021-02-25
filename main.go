package main

import (
	"log"
	
	"github.com/gin-gonic/gin"
	"github.com/nathan-tw/gin-rate-limiter/global"
	"github.com/nathan-tw/gin-rate-limiter/middleware"
	"github.com/nathan-tw/gin-rate-limiter/pkg/setting"
)

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Redis", &global.RedisSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	return nil
}

func init () {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}

func main() {
	
	r := gin.Default()
	r.Use(middleware.Limiter())
	r.Run()
}
