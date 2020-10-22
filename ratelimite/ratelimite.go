package ratelimite

import (
	"sync"
	"time"
)

//imple token bucket
//基于时间的令牌桶限流,尤其对于新增数据这种操作,理论上应该放在后端,目前直接对外暴露必须要做限流措施,以防恶意海量DDOS导致数据库超出负荷
type TokenBucket struct {
	rate    int
	surplus int
	endTime time.Time
	sync.Mutex
}

func (tb *TokenBucket) RequestTokenTimeout(n int, duration time.Duration) bool {
	if duration < time.Second {
		duration = time.Second
	}
	t := time.NewTimer(duration)
	defer t.Stop()
	select {
	case n = <-tb.requestToken(n):
		tb.consume(n)
		return true
	case <-t.C:
		return false
	}
}

func (tb *TokenBucket) consume(n int) {
	tb.surplus -= n
	tb.endTime = time.Now()
}

func (tb *TokenBucket) requestToken(n int) <-chan int {
	result := make(chan int)
	go func() {
		tb.Lock()
		ticker := time.NewTicker(16 * time.Millisecond)
		defer func() {
			ticker.Stop()
			tb.Unlock()
			result <- n
		}()
		for n > tb.surplus {
			//计算应该可以下放的请求数
			incrment := int(time.Now().Sub(tb.endTime).Seconds()) * tb.rate
			tb.surplus = incrment
			<-ticker.C
		}
	}()
	return result
}

func (tb *TokenBucket) RequestToken(n int) {
	tb.consume(<-tb.requestToken(n))
}

//@rate 速率(秒)
func NewTokenBucket(rate int) *TokenBucket {
	tb := TokenBucket{
		rate:    rate,
		surplus: 0,
		endTime: time.Now(),
	}
	return &tb
}
