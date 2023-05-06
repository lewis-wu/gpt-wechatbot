package cache

import (
	"fmt"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/869413421/wechatbot/util"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

var goCache *gocache.Cache

const sharedValue = "SHARED_VALUE"

func init() {
	goCache = gocache.New(5*time.Minute, 5*time.Minute)
}

func GetChatHistory(key string) ([]*dto.Message, bool) {
	lock(key)
	defer unLock(key)
	lr, ok := goCache.Get(key)
	if !ok {
		return nil, false
	}
	loopArray := lr.(*util.LoopArray)
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
	lr, ok := goCache.Get(key)
	if ok {
		loopArray := lr.(*util.LoopArray)
		loopArray.Push(value)
		return
	}
	la := util.NewLoopArray(config.LoadConfig().ChatMaxContext)
	la.Push(value)
	ttl := time.Duration(config.LoadConfig().ChatTTLTime) * time.Minute
	goCache.Set(key, la, ttl)
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
	goCache.Set(key, sharedValue, ttl)
}
func GetImageVar(key string) bool {
	_, ok := goCache.Get(key)
	if ok {
		goCache.Delete(key)
	}
	return ok
}
func BuildImageVarCacheKey(userName string, groupId string, isGroup bool) string {
	if isGroup {
		return fmt.Sprintf("image_var:room:%s:%s", groupId, userName)
	} else {
		return fmt.Sprintf("image_var:single:%s", userName)
	}
}
