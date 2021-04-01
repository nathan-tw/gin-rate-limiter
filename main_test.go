package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var (
	s      *miniredis.Miniredis
	testRS RedisServer
	r      *gin.Engine
	err    error
)

func setup() {
	gin.SetMode(gin.TestMode)
	s, _ = miniredis.Run()
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	testRS = RedisServer{
		timeout: 3600 * time.Second,
		maxIP:   10,
		ctx:     context.Background(),
		client:  client,
	}
	testRS.FlushAll()
	r = setupRouter(&testRS)
}

func tearDown() {
	s = nil
	testRS = RedisServer{}
	r = nil
	err = nil
}

func TestStatusCode(t *testing.T) {
	setup()
	for i := 0; i < testRS.maxIP+2; i++ {
		w := getResponse(r)
		if i < testRS.maxIP {
			assert.Equal(t, 200, w.Code)
		} else {
			assert.Equal(t, 429, w.Code)
		}
	}
	tearDown()
}

func TestRateLimitRemaining(t *testing.T) {
	setup()
	for i := 0; i < testRS.maxIP; i++ {
		w := getResponse(r)
		count, _ := strconv.Atoi(w.Result().Header.Get("X-RateLimit-Remaining"))
		assert.Equal(t, i, (testRS.maxIP-1)-count)
	}
	tearDown()
}

func TestRaceCondition(t *testing.T) {
	setup()
	var wg sync.WaitGroup
	wg.Add(testRS.maxIP)
	for i := 0; i < testRS.maxIP; i++ {
		go func(r *gin.Engine) {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			r.ServeHTTP(w, req)
			t.Log(w.Result().Header.Get("X-RateLimit-Remaining"))
		}(r)
	}
	wg.Wait()
	w := getResponse(r)
	assert.Equal(t, 429, w.Code)
	tearDown()
}

func getResponse(r *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	return w
}


