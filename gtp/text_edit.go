package gtp

import (
	"bytes"
	"encoding/json"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/869413421/wechatbot/util"
	"io"
	"log"
	"net/http"
	"strings"
)

func TextEdit(question string, userName string, groupId string, isGroup bool) (string, error) {
	validQuestion := strings.Replace(question, config.LoadConfig().TextEditKeyword, "", 1)
	inputInstructSeparatorIndex := strings.Index(validQuestion, config.LoadConfig().TextEditSeparator)
	var input string
	var instruction string
	if inputInstructSeparatorIndex == -1 {
		input = ""
		instruction = validQuestion
	} else {
		input = validQuestion[:inputInstructSeparatorIndex]
		instruction = validQuestion[(len(config.LoadConfig().TextEditSeparator) + inputInstructSeparatorIndex):]
	}

	reqBody := dto.TextEditReq{
		Model:       "text-davinci-edit-001",
		Input:       input,
		Instruction: instruction,
		Temperature: 0.7,
		TopP:        1,
		N:           1,
	}

	requestData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	log.Printf("GPT textEdit request text:%v", string(requestData))
	req, err := http.NewRequest("POST", config.BASEURL+"edits", bytes.NewBuffer(requestData))
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
		log.Printf("GPT textEdit response error: %s \n", string(body))
		return "ChatGPT响应错误", nil
	}
	textEditResp := &dto.TextEditResp{}
	err = json.Unmarshal(body, textEditResp)
	if err != nil {
		return "", err
	}
	var reply = ""
	for _, choice := range textEditResp.Choices {
		if len(choice.Text) > 0 {
			reply += choice.Text
		}
	}

	log.Printf("GPT textEdit response text: %s \n", reply)
	return reply, nil

}
