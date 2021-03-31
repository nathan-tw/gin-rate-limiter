package main

import (
	"context"
	"fmt"
	"strconv"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Database interface {
	GetKey(string) (string, string, error)
	SetKey(string) error
	CheckExist(string) bool
	IncreaseVisit(string) error
	FlushAll() error
	Lock() error
	Unlock() error
}

type RedisServer struct {
	timeout time.Duration
	client  *redis.Client
	ctx context.Context
	maxIP   int
}

func NewRedis(maxIP, timeout int, host, port string) (*RedisServer, error) {
	rs := new(RedisServer)
	addr := fmt.Sprintf("%v:%v", host, port)
	rs.client = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	err := rs.client.Ping(rs.client.Context()).Err()
	if err != nil {
		return nil, err
	}
	rs.ctx = context.Background()
	rs.timeout = time.Duration(timeout) * time.Second
	rs.maxIP = maxIP

	return rs, nil
}

func (rs *RedisServer) Set(ip string) error {
	err := rs.client.Set(rs.ctx, ip, "1", rs.timeout).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisServer) CheckExist(ip string) bool {
	_, err := rs.client.Get(rs.ctx, ip).Result()
	return err != redis.Nil
}

func (rs *RedisServer) Get(ip string) (string, string, error) {
	count, err := rs.client.Get(rs.client.Context(), ip).Int()
	if err != nil {
		return "", "", err
	}
	ttl := rs.client.TTL(rs.client.Context(), ip).Val().String()
	remaining := rs.maxIP - count

	return strconv.Itoa(remaining), ttl, nil
}

func (rs *RedisServer) IncreaseVisit(ip string) error {
	err := rs.client.Incr(rs.ctx, ip).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisServer) FlushAll() error {
	_, err := rs.client.FlushAll(rs.ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisServer) Lock() error {
	for {
		lockTimeout := time.Now().Add(10 * time.Second).Unix()
		ok, err := rs.client.SetNX(rs.ctx, "lock", lockTimeout, 0).Result()
		if err != nil {
			return err
		}

		// successfully get lock
		if ok {
			break
		}

		// the lock is taken
		TTL, err := rs.client.Get(rs.ctx, "lock").Int64()
		if err != nil {
			return err
		}

		// if expired
		curr := time.Now().Unix()
		if TTL <= curr {
			lockTimeout = time.Now().Add(10 * time.Second).Unix()

			// if multiple clients compete to snatch the lock
			response, err := rs.client.GetSet(rs.ctx, "lock", lockTimeout).Int64()
			if err != nil {
				return err
			}
			// response should not timeout if another client get the lock
			if curr > response {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func (rs *RedisServer) Unlock() error {
	now := time.Now().Unix()
	response, err := rs.client.Get(rs.ctx, "lock").Int64()
	if err != nil {
		return err
	}

	// shouldn't unlock due to successfully lock
	if now > response {
		return nil
	}

	response, err = rs.client.Del(rs.ctx, "lock").Result()
	if err != nil {
		return err
	}

	if response != 1 {
		return errors.New("Fail to delete lock")
	}

	return nil
}