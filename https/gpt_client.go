package https

import (
	"github.com/869413421/wechatbot/config"
	"net/http"
	"net/url"
	"time"
)

const DEFAULT_TIME_OUT = 60

func AddHeaderFroGpt(req *http.Request) {
	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
}

func GetGptClient() *http.Client {
	proxy := config.LoadConfig().Proxy
	timeOut := config.LoadConfig().GptTimeOut
	if timeOut <= 0 {
		timeOut = DEFAULT_TIME_OUT
	}
	var client *http.Client
	if len(proxy) == 0 {
		client = &http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		}
	} else {
		proxyAddr, _ := url.Parse(proxy)
		client = &http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyAddr),
			},
		}
	}
	return client
}
