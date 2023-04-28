package handlers

import (
	"github.com/869413421/wechatbot/pool"
	"github.com/eatmoreapple/openwechat"
	"log"
)

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	if msg.IsText() {
		log.Printf("handler received text msg => %v", msg.Content)
	} else {
		log.Printf("handler received non text msg, msg type => %v", msg.MsgType)
	}

	if err := pool.GetPool().Submit(func() { doHandle(msg) }); err != nil {
		log.Printf("submit to ants pool error => %v", err)
	}
}
