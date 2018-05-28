package tokenbucket

import (
	"errors"
	"sync"
)

// Queue : ring queue struct
type Queue struct {
	mtx      sync.RWMutex
	head     int64
	tail     int64
	capacity int64
	elemcnt  int64
	data     []interface{}
}

// NewQueue returns a new Queue instance or an error
func NewQueue(cap int64) (*Queue, error) {
	if cap <= 0 {
		return nil, errors.New("invalid capacity specified")
	}
	return &Queue{
		capacity: cap,
		data:     make([]interface{}, cap),
	}, nil
}

// IsEmpty checks if queue is empty by element count
func (q *Queue) IsEmpty() bool {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return q.elemcnt == 0
}

// IsFull checks if queue is full by element count
func (q *Queue) IsFull() bool {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return q.elemcnt >= q.capacity
}

// Available checks if there's at least n block is not-used
func (q *Queue) Available(n int64) bool {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return q.capacity-q.elemcnt >= n
}

// Len returns a element count of the queue
func (q *Queue) Len() int64 {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return q.elemcnt
}

// LastEnqueue returns the lastest element in queue
func (q *Queue) LastEnqueue() interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.elemcnt == 0 {
		return nil
	}
	if q.head > 0 {
		return q.data[q.head-1]
	}
	return q.data[q.capacity-1]
}

// FirstEnqueue returns the oldest element in queue
func (q *Queue) FirstEnqueue() interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.elemcnt == 0 {
		return nil
	}
	return q.data[q.tail]
}

// Data returns all elements in queue
func (q *Queue) Data() []interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.elemcnt == 0 {
		return []interface{}{}
	}
	ret := make([]interface{}, q.elemcnt)
	for i := int64(0); i < q.elemcnt; i++ {
		idx := q.tail + i
		if idx >= q.capacity {
			idx -= q.capacity
		}
		ret[i] = q.data[idx]
	}
	return ret
}

// EnqueueOne try to enqueue an element
func (q *Queue) EnqueueOne(d interface{}) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	// update tail index
	if q.head == q.tail && q.elemcnt >= q.capacity {
		q.tail++
		if q.tail >= q.capacity {
			q.tail = 0
		}
	}
	// push data on head index
	q.data[q.head] = d
	// update head index
	q.head++
	if q.head >= q.capacity {
		q.head = 0
	}
	// update element count
	if q.elemcnt < q.capacity {
		q.elemcnt++
	}
}

// AtomEnqueue enqueue an element for n times
func (q *Queue) AtomEnqueue(d interface{}, n int64) bool {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.capacity-q.elemcnt < n {
		return false
	}
	q.enqueue(d, n)
	return true
}

// EnqueueN enqueue a list of elements
func (q *Queue) EnqueueN(data []interface{}) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	for _, d := range data {
		if q.head == q.tail && q.elemcnt >= q.capacity {
			q.tail++
			if q.tail >= q.capacity {
				q.tail = 0
			}
		}
		q.data[q.head] = d
		q.head++
		if q.head >= q.capacity {
			q.head = 0
		}
		if q.elemcnt < q.capacity {
			q.elemcnt++
		}
	}
}

// DequeueN dequeue a number of elements
func (q *Queue) DequeueN(n int64) []interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	retlen := n
	if retlen > q.elemcnt {
		retlen = q.elemcnt
	}
	ret := make([]interface{}, retlen)
	for i := int64(0); i < retlen; i++ {
		ret[i] = q.data[q.tail]
		q.tail++
		if q.tail >= q.capacity {
			q.tail = 0
		}
	}
	q.elemcnt -= retlen
	return ret
}

// DequeueBy dequeue elements using a critier-function
func (q *Queue) DequeueBy(fc func(interface{}) bool) []interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.elemcnt == 0 {
		return []interface{}{}
	}

	ret := []interface{}{}
	for i := int64(0); i < q.elemcnt; i++ {
		if fc(q.data[q.tail]) == false {
			break
		}
		ret = append(ret, q.data[q.tail])
		q.tail++
		if q.tail >= q.capacity {
			q.tail = 0
		}
	}
	q.elemcnt -= int64(len(ret))
	return ret
}

// enqueue an element for n times
func (q *Queue) enqueue(d interface{}, n int64) {
	if n > q.capacity {
		n = q.capacity
	}

	for i := int64(0); i < n; i++ {
		if q.head == q.tail && q.elemcnt >= q.capacity {
			q.tail++
			if q.tail >= q.capacity {
				q.tail = 0
			}
		}
		q.data[q.head] = d
		q.head++
		if q.head >= q.capacity {
			q.head = 0
		}
		if q.elemcnt < q.capacity {
			q.elemcnt++
		}
	}
}
