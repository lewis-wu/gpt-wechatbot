package cache

import (
	"sync"
	"sync/atomic"
)

var lockMap = sync.Map{}

func lock(key string) {
	var mutex int32 = 0
	m, _ := lockMap.LoadOrStore(key, &mutex)
	m32 := m.(*int32)
	for !atomic.CompareAndSwapInt32(m32, 0, 1) {
		return
	}
}
func unLock(key string) {
	m, ok := lockMap.Load(key)
	if !ok {
		panic("Lock has released before.")
	}
	atomic.StoreInt32(m.(*int32), 0)
}
