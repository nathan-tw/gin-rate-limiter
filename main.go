package main

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setupRouter(db *RedisServer) *gin.Engine {
	r := gin.Default()
	r.GET("/", Limiter(db))
	return r
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)

	config, err := initConfig()
	if err != nil {
		log.Fatalf("get config err: %v", err)
	}
	config.UnmarshalKey("Redis", &redisSetting)
	config.UnmarshalKey("App", &appSetting)
}

func main() {
	db, err := NewRedis(appSetting.MaxIP, appSetting.Timeout, redisSetting.Host, redisSetting.Port)
	if err != nil {
		log.Fatalf("error while new redis server: %v", err)
	}
	r := setupRouter(db)
	r.Run()
}
