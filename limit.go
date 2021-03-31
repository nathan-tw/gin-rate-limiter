package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)



func Limiter (rs *RedisServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if err := rs.Lock(); err != nil {
			log.Fatalf("database server lock error: %v", err)
		}

		if !rs.CheckExist(ip) {
			err := rs.Set(ip)
			if err != nil {
				log.Fatalf("middleware.Limiter err: %v", err)
			}
		}

		XRemaining, XReset, err := rs.Get(ip)
		if err != nil {
			log.Fatalf("middleware.Limiter err: %v", err)
		}
		XRemaingInt, err := strconv.Atoi(XRemaining)
		if err != nil {
			log.Fatalf("middleware.Limiter err: %v", err)
		}
		if  XRemaingInt < 0 {
			c.JSON(429, nil)
		} else {
			rs.IncreaseVisit(ip)
			c.Header("X-RateLimit-Remaining", XRemaining)
			c.Header("X-RateLimit-Reset", XReset )
			c.JSON(200, nil)
		}
		if err := rs.Unlock(); err != nil {
			log.Fatalf("redis server unlock error: %v", err)
		}

	}
}