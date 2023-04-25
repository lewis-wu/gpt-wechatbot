package limit

import (
	"github.com/869413421/wechatbot/config"
	"golang.org/x/time/rate"
	"sync"
	"sync/atomic"
)

var limiterMap *sync.Map

func init() {
	limiterMap = &sync.Map{}
}

type limiterHolder struct {
	limiter *rate.Limiter
	mutex   int32
}

func (holder *limiterHolder) get() *rate.Limiter {
	for {
		if holder.limiter != nil {
			return holder.limiter
		}
		if atomic.CompareAndSwapInt32(&holder.mutex, 0, 1) {
			limitPerMinute := config.LoadConfig().GptLimitPerMinute
			lim := rate.NewLimiter(rate.Limit(float64(limitPerMinute)/60), 10)
			holder.limiter = lim
			return lim
		}
	}
}

func ShouldLimit(key string) bool {
	limHolder := &limiterHolder{
		limiter: nil,
		mutex:   0,
	}
	limH, _ := limiterMap.LoadOrStore(key, limHolder)
	limiter := limH.(*limiterHolder).get()
	return !limiter.Allow()
}
