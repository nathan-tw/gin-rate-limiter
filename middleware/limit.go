package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nathan-tw/gin-rate-limiter/global"
	"github.com/nathan-tw/gin-rate-limiter/redis"
)



func Limiter () gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		// if ip does not exist, set to 1
		if !redis.RedisServer.CheckExist(ip) {
			err := redis.RedisServer.SetKey(ip)
			if err != nil {
				global.Logger.Fatalf("middleware.Limiter err: %v", err)
			}
		}
		XRemaining, XReset, err := redis.RedisServer.GetKey(ip)
		if err != nil {
			global.Logger.Fatalf("middleware.Limiter err: %v", err)
		}
		XRemaingInt, err := strconv.Atoi(XRemaining)
		if err != nil {
			global.Logger.Fatalf("middleware.Limiter err: %v", err)
		}
		if  XRemaingInt < 0 {
			c.JSON(429, nil)
		} else {
			redis.RedisServer.IncreaseVisit(ip)
			c.Header("X-RateLimit-Remaining", XRemaining)
			c.Header("X-RateLimit-Reset", XReset )
			c.JSON(200, nil)
		}

	}
}