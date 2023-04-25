package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

const BASEURL = "https://api.openai.com/v1/"
const defaultImageKeyword = "[图片]"
const defaultTextEditKeyword = "[TE]"
const defaultTextEditSeparator = "[TES]"
const defaultGptLimitPerMinute = 3

// Configuration 项目配置
type Configuration struct {
	ApiKey               string `json:"api_key"`                // gtp apikey
	AutoPass             bool   `json:"auto_pass"`              // 自动通过好友
	Proxy                string `json:"proxy"`                  //代理 http(s)://xxx.xxx:port
	ChatMaxContext       int    `json:"chat_max_context"`       //保存的最大聊天上下文记录数
	ChatTTLTime          int    `json:"chat_ttl_time"`          //聊天上下文保存的时间（小时）
	GptTimeOut           int    `json:"gpt_time_out"`           //gpt接口超时时间（秒）
	GenerateImageKeyword string `json:"generate_image_keyword"` //生成图片时所需的聊天关键词
	TextEditKeyword      string `json:"text_edit_keyword"`      //文本编辑的关键词
	TextEditSeparator    string `json:"text_edit_separator"`    //文本编辑的问题内容/问题建议分隔符
	GptLimitPerMinute    int    `json:"gpt_limit_per_minute"`   //chat请求限速（每分钟请求次数上限）
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
		if len(strings.TrimSpace(config.GenerateImageKeyword)) == 0 {
			config.GenerateImageKeyword = defaultImageKeyword
		}
		if len(strings.TrimSpace(config.TextEditKeyword)) == 0 {
			config.TextEditKeyword = defaultTextEditKeyword
		}
		if len(strings.TrimSpace(config.TextEditSeparator)) == 0 {
			config.TextEditSeparator = defaultTextEditSeparator
		}
		if config.GptLimitPerMinute <= 0 {
			config.GptLimitPerMinute = defaultGptLimitPerMinute
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
