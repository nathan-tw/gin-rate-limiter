package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nathan-tw/gin-rate-limiter/global"
)

type Database interface {
	GetKey(string) (string, string, error)
	SetKey(string) error
	CheckExist(string) bool
	IncreaseVisit(string) error
}

type Server struct {
	timeout time.Duration
	client  *redis.Client
	maxIP   int
}

func NewRedis(maxIP, timeout int) *Server {
	db := new(Server)

	addr := fmt.Sprintf("%v:%v", global.RedisSetting.Host, global.RedisSetting.Port)
	db.client = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	db.timeout = time.Duration(timeout) * time.Second
	db.maxIP = maxIP

	return db
}

func (d *Server) SetKey(ip string,) error {
	err := d.client.Set(d.client.Context(), ip, 1, d.timeout).Err()
	if err != nil {
		return err
	}
	return nil
}


func (d *Server) GetKey(ip string) (string, string, error) {
	count, err := d.client.Get(d.client.Context(), ip).Int()
	if err != nil {
		return "", "", err
	}
	ttl := d.client.TTL(d.client.Context(), ip).Val().String()
	remaining := d.maxIP - count

	return strconv.Itoa(remaining), ttl, nil
}

func (d *Server) CheckExist(ip string) bool {
	_, err := d.client.Get(d.client.Context(), ip).Result()
	if err == redis.Nil {
		return false
	}
	return true
}

func (d *Server) IncreaseVisit(ip string) error {
	err := d.client.Incr(d.client.Context(), ip).Err()
	if err != nil {
		return err
	}
	return nil
}
