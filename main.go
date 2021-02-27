package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nathan-tw/gin-rate-limiter/global"
	"github.com/nathan-tw/gin-rate-limiter/middleware"
	"github.com/nathan-tw/gin-rate-limiter/pkg/logger"
	"github.com/nathan-tw/gin-rate-limiter/pkg/setting"
	"github.com/nathan-tw/gin-rate-limiter/redis"
	"gopkg.in/natefinch/lumberjack.v2"
)



func init () {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err %v", err)
	}	
	err = setupRedis()
	if err != nil {
		log.Fatalf("init.setupRedis err: %v", err)
	}
	fmt.Println(redis.RedisServer.CheckExist("1"))
}

func main() {
	global.Logger.Info("Start using gin-rate-limiter")
	r := gin.Default()
	r.Use(middleware.Limiter())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"foo": "bar",
		})
	})
	r.Run()
}

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
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFlieExt
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:	fileName,
		MaxSize:	600,
		MaxAge: 10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func setupRedis() error {
	const (
		maxIP = 1000
		timeout = 3600
	)
	db, err := redis.NewRedis(maxIP, timeout)
	redis.RedisServer = db
	if err != nil {
		return err
	}
	return nil
}