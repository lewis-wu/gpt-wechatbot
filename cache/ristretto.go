package cache

import (
	"fmt"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/dgraph-io/ristretto"
	"time"
)

const SingleKeyCost = 1

var cache *ristretto.Cache

const sharedValue = "SHARED_VALUE"

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
	la := NewLoopArray(config.LoadConfig().ChatMaxContext)
	la.Push(value)
	ttl := time.Duration(config.LoadConfig().ChatTTLTime) * time.Minute
	cache.SetWithTTL(key, la, SingleKeyCost, ttl)
}

func BuildChatHistoryCacheKey(userName string, groupId string, isGroup bool) string {
	if isGroup {
		return fmt.Sprintf("room:%s:%s", groupId, userName)
	} else {
		return fmt.Sprintf("single:%s", userName)
	}
}

func AddImageVar(key string) {
	ttl := time.Duration(config.LoadConfig().ImageVariationChatTTL) * time.Second
	cache.SetWithTTL(key, sharedValue, SingleKeyCost, ttl)
}
func GetImageVar(key string) bool {
	_, ok := cache.Get(key)
	cache.Del(key)
	return ok
}
func BuildImageVarCacheKey(userName string, groupId string, isGroup bool) string {
	if isGroup {
		return fmt.Sprintf("image_var:room:%s:%s", groupId, userName)
	} else {
		return fmt.Sprintf("image_var:single:%s", userName)
	}
}
