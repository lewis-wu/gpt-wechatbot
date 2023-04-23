package gtp

import (
	"bytes"
	"encoding/json"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/869413421/wechatbot/https"
	"io"
	"log"
	"net/http"
)

func TextEdit(question string, userName string, groupId string, isGroup bool) (string, error) {
	reqBody := dto.TextEditReq{
		Model:       "gpt-3.5-turbo",
		Input:       question,
		Instruction: "",
		Temperature: 0.7,
		TopP:        1,
		N:           1,
	}

	requestData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	log.Printf("GPT textComplete request text:%v", string(requestData))
	req, err := http.NewRequest("POST", config.BASEURL+"edits", bytes.NewBuffer(requestData))
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
	textCompleteResp := &dto.TextEditResp{}
	err = json.Unmarshal(body, textCompleteResp)
	if err != nil {
		return "", err
	}
	var reply = ""
	for _, choice := range textCompleteResp.Choices {
		if len(choice.Text) > 0 {
			reply += choice.Text
		}
	}

	log.Printf("GPT textComplete response text: %s \n", reply)
	return reply, nil

}
