package gtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/869413421/wechatbot/cache"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/869413421/wechatbot/https"
	"io"
	"log"
	"net/http"
	"strings"
)

func ChatCompletions(question string, userName string, groupId string, isGroup bool) (string, error) {
	messages := make([]*dto.Message, 0, 10)
	key := buildCacheKey(userName, groupId, isGroup)
	chatHistory, ok := cache.GetChatHistory(key)
	if ok {
		for _, s := range chatHistory {
			messages = append(messages, s)
		}
	}
	curMessage := &dto.Message{
		Role:    "user",
		Content: question,
	}
	messages = append(messages, curMessage)
	reqBody := dto.ChatCompleteReq{
		Model:            "gpt-3.5-turbo",
		Messages:         messages,
		MaxTokens:        3000,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Stream:           false,
	}

	requestData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	log.Printf("GPT chatComplete request text:%v", string(requestData))
	req, err := http.NewRequest("POST", config.BASEURL+"chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}
	https.AddHeaderForGpt(req)
	client := https.GetGptClient()

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	chatCompleteResp := &dto.ChatCompleteResp{}
	err = json.Unmarshal(body, chatCompleteResp)
	if err != nil {
		return "", err
	}
	//成功获取结果才将question放到上下文中
	cache.AddChatHistory(key, curMessage)
	var reply = ""
	for _, choice := range chatCompleteResp.Choices {
		if len(strings.TrimSpace(choice.Message.Content)) > 0 {
			reply += choice.Message.Content
			cache.AddChatHistory(key, choice.Message)
		}
	}
	log.Printf("GPT chatComplete response text: %s \n", reply)
	return reply, nil

}
func buildCacheKey(userName string, groupId string, isGroup bool) string {
	if isGroup {
		return fmt.Sprintf("room:%s:%s", groupId, userName)
	} else {
		return fmt.Sprintf("single:%s", userName)
	}
}
