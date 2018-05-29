# token-bucket

[![Build Status](https://api.travis-ci.org/moxiaomomo/token-bucket.svg)](https://travis-ci.org/moxiaomomo/token-bucket)
[![Go Report Card](https://goreportcard.com/badge/github.com/moxiaomomo/token-bucket)](https://goreportcard.com/report/github.com/moxiaomomo/token-bucket)
结合环形队列(ring queue)和令牌桶(token bucket)算法实现的一个限流模块。

### 测试用例

```
$ go test
[1 2 3 4 5 6]
[1 2 3 4 5 6 7 8]
[3 4 5 6 7 8 9 10]
[9 10 3 4 5 6 7 8]
[10 11 12 13 14 15 16 17]
[10 11 12 13 14]
[15 16 17]
[]
[9 10 haha 15 3 6]
[9 10]
take-n:1 time-used:3.327341023s suc-cnt:4000 fail-cnt:9996000
take-n:2 time-used:3.334162956s suc-cnt:2000 fail-cnt:9998000
take-n:1 time-used:1.257582083s suc-cnt:1499 fail-cnt:1001
take-n:1 time-used:2.250956732s suc-cnt:2500 fail-cnt:0
take-n:2 time-used:1.254886958s suc-cnt:999 fail-cnt:1501
take-n:2 time-used:4.251330391s suc-cnt:2499 fail-cnt:1
take-n:2 time-used:2.251117604s suc-cnt:2500 fail-cnt:0
PASS
ok      github.com/moxiaomomo/token-bucket23.931s
```
