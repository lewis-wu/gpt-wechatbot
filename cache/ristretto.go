package cache

import (
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/dgraph-io/ristretto"
	"time"
)

const SingleKeyCost = 1
const DefaultTTLTime = 1 * time.Hour
const DefaultChatMaxContext = 2

var cache *ristretto.Cache
var chatMaxContext int
var chatTTLTime time.Duration

func init() {
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // Num keys to track frequency of (10M).
		MaxCost:     1 << 30, // Maximum cost of cache (1GB).
		BufferItems: 64,      // Number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	num := config.LoadConfig().ChatMaxContext
	if num <= 0 {
		chatMaxContext = DefaultChatMaxContext
	} else {
		chatMaxContext = num
	}
	ttlTime := config.LoadConfig().ChatTTLTime
	if ttlTime <= 0 {
		chatTTLTime = DefaultTTLTime
	} else {
		chatTTLTime = time.Duration(ttlTime) * time.Hour
	}
}

func GetChatHistory(key string) ([]*dto.Message, bool) {
	lock(key)
	defer unLock(key)
	lr, ok := cache.Get(key)
	if !ok {
		return nil, false
	}
	loopArray := lr.(*LoopArray)
	data := loopArray.Clone()
	result := make([]*dto.Message, 0, len(data))
	for _, v := range data {
		if v != nil {
			result = append(result, v.(*dto.Message))
		}
	}
	return result, true
}
func AddChatHistory(key string, value *dto.Message) {
	lock(key)
	defer unLock(key)
	lr, ok := cache.Get(key)
	if ok {
		loopArray := lr.(*LoopArray)
		loopArray.Push(value)
		return
	}
	la := NewLoopArray(chatMaxContext)
	la.Push(value)
	cache.SetWithTTL(key, la, SingleKeyCost, chatTTLTime)
}
