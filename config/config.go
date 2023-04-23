package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

const BASEURL = "https://api.openai.com/v1/"

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 自动通过好友
	AutoPass       bool   `json:"auto_pass"`
	Proxy          string `json:"proxy"`
	ChatMaxContext int    `json:"chat_max_context"` //保存的最大聊天上下文记录数
	ChatTTLTime    int    `json:"chat_ttl_time"`    //聊天上下文保存的时间（小时）
	GptTimeOut     int    `json:"gpt_time_out"`     //gpt接口超时时间（秒）
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
			return
		}
		//// 如果环境变量有配置，读取环境变量
		//ApiKey := os.Getenv("ApiKey")
		//AutoPass := os.Getenv("AutoPass")
		//if ApiKey != "" {
		//	config.ApiKey = ApiKey
		//}
		//if AutoPass == "true" {
		//	config.AutoPass = true
		//}
	})
	return config
}
