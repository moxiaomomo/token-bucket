package tokenbucket

import (
	"fmt"
	"testing"
	"time"
)

func test_Take(tb *TokenBucket, n int64) {
	st := time.Now()
	suc, fail := 0, 0
	for i := 0; i < 10000000; i++ {
		res, _ := tb.Take(n)
		if res {
			suc++
		} else {
			fail++
		}
	}
	fmt.Printf("take-n:%d time-used:%s suc-cnt:%d fail-cnt:%d\n",
		n, time.Now().Sub(st).String(), suc, fail)
}

func test_Wait(tb *TokenBucket, wait time.Duration, n int64) {
	st := time.Now()
	suc, fail := 0, 0
	ch := make(chan int, 100)
	for _ = range []int{1, 2, 3, 4, 5} {
		for i := 0; i < 500; i++ {
			go func(ch chan int) {
				res := tb.Wait(n, wait)
				if res {
					ch <- 1
				} else {
					ch <- 0
				}
			}(ch)
		}
		time.Sleep(250 * time.Millisecond)
	}

	for i := 0; i < 2500; i++ {
		res := <-ch
		if res == 1 {
			suc++
		} else {
			fail++
		}
	}
	fmt.Printf("take-n:%d time-used:%s suc-cnt:%d fail-cnt:%d\n",
		n, time.Now().Sub(st).String(), suc, fail)
}

func Test_tokenbucket(t *testing.T) {
	// limit to 1000 per second
	tb, err := NewTokenBucket(1000, time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	test_Take(tb, 1)

	time.Sleep(time.Second)
	test_Take(tb, 2)

	time.Sleep(time.Second)
	test_Wait(tb, 100*time.Millisecond, 1)

	time.Sleep(time.Second)
	test_Wait(tb, 10*time.Second, 1)

	time.Sleep(time.Second)
	test_Wait(tb, 100*time.Millisecond, 2)

	time.Sleep(time.Second)
	test_Wait(tb, 10*time.Second, 2)

	// limit to 2000 per second
	tb, err = NewTokenBucket(2000, time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second)
	test_Wait(tb, 10*time.Second, 2)
}
