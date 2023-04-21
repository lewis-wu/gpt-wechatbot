package cache

import "sync/atomic"

// LoopArray is a fixed capacity concurrent safe array
type LoopArray struct {
	data  []interface{}
	cap   int
	tail  int32
	head  int32
	n     int32
	mutex int32
}

// NewLoopArray creates a fixed capacity concurrent safe array
func NewLoopArray(cap int) *LoopArray {
	return &LoopArray{
		data:  make([]interface{}, cap),
		cap:   cap,
		tail:  0,
		head:  0,
		n:     0,
		mutex: 0,
	}
}

// Push adds an element to the array
func (l *LoopArray) Push(x interface{}) {
	if !atomic.CompareAndSwapInt32(&l.mutex, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&l.mutex, 0)

	l.data[l.tail] = x
	l.tail++
	if l.tail == int32(l.cap) {
		l.tail = 0
	}
	if l.n < int32(l.cap) {
		l.n++
	} else {
		l.head = l.tail
	}
}

// Get returns the element at given index
func (l *LoopArray) Get(i int32) interface{} {
	if atomic.LoadInt32(&l.n) == 0 || i >= atomic.LoadInt32(&l.n) {
		return nil
	}
	return l.data[(l.head+i)%int32(l.cap)]
}

// Clone returns a copy of the array
func (l *LoopArray) Clone() []interface{} {
	if atomic.LoadInt32(&l.n) == 0 {
		return nil
	}
	result := make([]interface{}, 0, atomic.LoadInt32(&l.n))
	for i := int32(0); i < atomic.LoadInt32(&l.n); i++ {
		result = append(result, l.Get(i))
	}
	return result
}
