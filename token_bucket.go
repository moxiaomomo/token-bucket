package tokenbucket

import (
	"sync"
	"time"
)

// token bucket struct for current-limiting
type TokenBucket struct {
	mutex    sync.Mutex
	capacity int64
	interval time.Duration
	tsQ      *Queue
}

var count = 0

func NewTokenBucket(limit int64, interval time.Duration) (*TokenBucket, error) {
	rqueue, err := NewQueue(limit)
	if err != nil {
		return nil, err
	}
	tb := &TokenBucket{
		interval: interval,
		capacity: limit,
		tsQ:      rqueue,
	}
	return tb, nil
}

// try to take n token-buckets
// return (true, waittime=0) if not limit reaches,
// otherwise return (false, waittime>0)
func (tb *TokenBucket) Take(n int64) (bool, time.Duration) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	if tb.tsQ.AtomEnqueue(now, n) {
		return true, 0
	}

	_ = tb.tsQ.DequeueBy(func(val interface{}) bool {
		if v, ok := val.(time.Time); ok {
			if now.Sub(v) >= tb.interval {
				return true
			}
		}
		return false
	})

	oldest, _ := tb.tsQ.FirstEnqueue().(time.Time)
	return false, tb.interval - now.Sub(oldest)
}

// wait for token-buckets until timeout,
// return true if take buckets succeeded,
// else return false
func (tb *TokenBucket) Wait(n int64, t time.Duration) bool {
	timeout := t
	for {
		suc, waitT := tb.Take(n)
		if suc {
			return true
		}

		timeout -= waitT
		if timeout > 0 {
			time.Sleep(waitT)
		} else {
			return false
		}
	}
}
