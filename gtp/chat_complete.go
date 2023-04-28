package gtp

import (
	"bytes"
	"encoding/json"
	"github.com/869413421/wechatbot/cache"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/869413421/wechatbot/util"
	"io"
	"log"
	"net/http"
	"strings"
)

const maxTokens = 4000

func ChatCompletions(question string, userName string, groupId string, isGroup bool) (string, error) {
	historyMaxToken := maxTokens - len(question)
	if historyMaxToken < 0 {
		return "提问文本超过了最大字数，无法回答", nil
	}
	messages := make([]*dto.Message, 0, 10)
	key := cache.BuildChatHistoryCacheKey(userName, groupId, isGroup)
	chatHistory, ok := cache.GetChatHistory(key)
	if ok {
		messages = buildUseChatHistory(chatHistory, historyMaxToken, messages)
	}
	curMessage := &dto.Message{
		Role:    "user",
		Content: question,
	}
	messages = append(messages, curMessage)
	reqBody := dto.ChatCompleteReq{
		Model:            "gpt-3.5-turbo",
		Messages:         messages,
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
	util.AddHeaderForGpt(req)
	client := util.GetGptClient()

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		log.Printf("GPT chatComplete response error: %s \n", string(body))
		return "ChatGPT响应错误", nil
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

func buildUseChatHistory(chatHistory []*dto.Message, historyMaxToken int, messages []*dto.Message) []*dto.Message {
	if len(chatHistory) == 0 {
		return messages
	}
	historyTokenCount := 0
	usedHistory := make([]*dto.Message, 0, len(chatHistory))
	for i := len(chatHistory) - 1; i >= 0; i-- {
		message := chatHistory[i]
		historyTokenCount += len(message.Content)
		if historyMaxToken-historyTokenCount < 0 {
			break
		}
		usedHistory = append(usedHistory, message)
	}
	if len(usedHistory) == 0 {
		return messages
	}

	for i := len(usedHistory) - 1; i >= 0; i-- {
		messages = append(messages, usedHistory[i])
	}
	return messages
}
