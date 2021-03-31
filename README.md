# gin-rate-limiter

## Requirement

Dcard 每天午夜都有大量使用者湧入抽卡，為了不讓伺服器過載，請設計一個 middleware：

-   限制每小時來自同一個 IP 的請求數量不得超過 1000
-   在 response headers 中加入剩餘的請求數量 (X-RateLimit-Remaining) 以及 rate limit 歸零的時間 (X-RateLimit-Reset)
-   如果超過限制的話就回傳 429 (Too Many Requests)
-   可以使用各種資料庫達成

## Quick start

demo的程式將跑在localhost:8080

```bash
$ docker-compose up
```

## Race Condition

由於來自同一個ip的request有可能是非同步的，雖然redis的設計主要為單執行緒，
但因為一次的request將有數個transactions要做，有可能在get和set之間有其他request的get或set造成race condition，
因此需要設計一個lock保證資料完整性。我是參考[redis官網SETNX語法](https://redis.io/commands/setnx)補充的演算法，實現一個簡單的lock。
redis lock有很多實現方式，例如對多個redis做分散鎖的redlock或我較不熟悉的zookeeper，我會選這個方式除了比較簡單實現外，考量到來自同一個ip的request不應如此複雜所以選擇這個方式。

#### Algo Description

SETNX指令會試著將某個key設值，但如果該key-value pair已存在則會失敗返回0，也就是**SET** if **N**ot e**X**ist。
因此可以藉由這個指令設一個lock，每次request來時先對這個lock下setnx，將值設為expired time，如果成功設立代表該client取得lock，反之則未取得須等待一段時間重新嘗試。考慮到如果某個client因為一些因素無法釋放鎖。鎖過期時會發生client競爭鎖的情形，如果SETNX失敗了則先GET確認鎖是否過期，如果是的話則使用GETSET，接著將GETSET的結果再次確認是否過期，如此一來，先搶到鎖的人會率先因GETSET更改了expired time，使得後使用GETSET的client取得lock並未過期。
