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
	"strings"
)

func GenerateImage(prompt string, userName string, groupId string, isGroup bool) (string, error) {
	createImageReq := &dto.CreateImageReq{
		Prompt:         strings.Replace(prompt, config.LoadConfig().GenerateImageKeyword, "", 1),
		N:              1,
		Size:           dto.IMAGE_SIZE_1024,
		ResponseFormat: dto.IMAGE_FROMAT_URL,
	}
	requestData, err := json.Marshal(createImageReq)
	if err != nil {
		return "", err
	}
	log.Printf("GPT createImage request text:%v", string(requestData))
	req, err := http.NewRequest("POST", config.BASEURL+"images/generations", bytes.NewBuffer(requestData))
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
	genImageResp := &dto.ImageResp{}
	err = json.Unmarshal(body, genImageResp)
	if err != nil {
		return "", err
	}
	var imageURL = ""
	if genImageResp.Created == 0 {
		log.Printf("GPT createImage error:%v", string(body))
		return "可能已触发违禁词，不能正常生成图片", nil
	}
	for _, content := range genImageResp.ImageContents {
		if len(content.URL) > 0 {
			imageURL += content.URL + "\n"
		}
	}
	log.Printf("GPT createImage response text:%v", string(requestData))
	return imageURL, nil
}
