# gin-rate-limiter

## Requirement

Dcard 每天午夜都有大量使用者湧入抽卡，為了不讓伺服器過載，請設計一個 middleware：
* 限制每小時來自同一個 IP 的請求數量不得超過 1000
* 在 response headers 中加入剩餘的請求數量 (X-RateLimit-Remaining) 以及 rate limit 歸零的時間 (X-RateLimit-Reset)
* 如果超過限制的話就回傳 429 (Too Many Requests)
* 可以使用各種資料庫達成

## Quick Start

```go
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
```