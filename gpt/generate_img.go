package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/869413421/wechatbot/util"
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
		ResponseFormat: dto.IMAGE_FROMAT_BASE64,
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
		log.Printf("GPT createImage error:%v", string(body))
		return "", errors.New("create image failed")
	}
	genImageResp := &dto.ImageResp{}
	err = json.Unmarshal(body, genImageResp)
	if err != nil {
		return "", err
	}
	var imageURL = ""
	for _, content := range genImageResp.ImageContents {
		imageURL = content.B64Json
		break
	}
	return imageURL, nil
}
