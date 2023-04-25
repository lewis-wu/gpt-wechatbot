package https

import (
	"github.com/869413421/wechatbot/config"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

const DEFAULT_TIME_OUT = 60

var mutex int32 = 0
var client *http.Client = nil

func AddHeaderForGpt(req *http.Request) {
	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
}

func GetGptClient() *http.Client {
	for {
		if client != nil {
			return client
		}
		if atomic.CompareAndSwapInt32(&mutex, 0, 1) {
			proxy := config.LoadConfig().Proxy
			timeOut := config.LoadConfig().GptTimeOut
			if timeOut <= 0 {
				timeOut = DEFAULT_TIME_OUT
			}
			var tmpClient *http.Client
			if len(proxy) == 0 {
				tmpClient = &http.Client{
					Timeout: time.Duration(timeOut) * time.Second,
				}
			} else {
				proxyAddr, _ := url.Parse(proxy)
				tmpClient = &http.Client{
					Timeout: time.Duration(timeOut) * time.Second,
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyAddr),
					},
				}
			}
			client = tmpClient
			return client
		}
	}
}
