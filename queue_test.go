package tokenbucket

import (
	"fmt"
	"testing"
)

func compare(val interface{}) bool {
	if v, ok := val.(int); ok {
		if v > 5 {
			return true
		}
	}
	return false
}

func Test_queue(t *testing.T) {
	q, _ := NewQueue(8)
	q.EnqueueN([]interface{}{1, 2, 3, 4, 5, 6})
	fmt.Printf("%+v\n", q.Data())

	q.EnqueueN([]interface{}{7, 8})
	fmt.Printf("%+v\n", q.Data())

	q.EnqueueN([]interface{}{9, 10})
	fmt.Printf("%+v\n", q.Data())
	fmt.Printf("%+v\n", q.data)

	q.EnqueueN([]interface{}{9, 10, 11, 12, 13, 14, 15, 16, 17})
	fmt.Printf("%+v\n", q.Data())

	ret := q.DequeueN(5)
	fmt.Printf("%+v\n", ret)

	ret = q.DequeueN(5)
	fmt.Printf("%+v\n", ret)

	ret = q.DequeueN(5)
	fmt.Printf("%+v\n", ret)

	q.EnqueueN([]interface{}{9, 10, "haha", 15, 3, 6})
	fmt.Printf("%+v\n", q.Data())

	ret = q.DequeueBy(compare)
	fmt.Printf("%+v\n", ret)
}
