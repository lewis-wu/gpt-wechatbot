package gtp

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/869413421/wechatbot/util"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func ImageVariation(jpgImg image.Image, userName string, groupId string, isGroup bool) (string, error) {
	pngFile, err := util.Jpg2PngAndResize(jpgImg, 1024)
	if err != nil {
		return "", err
	}

	defer util.DeleteImage(pngFile)
	req, err := buildRequest(pngFile)
	if err != nil {
		return "", err
	}
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
		log.Printf("GPT ImageVariation error:%v\n", string(body))
		return "", errors.New("image variation failed")
	}
	variationImageResp := &dto.ImageResp{}
	err = json.Unmarshal(body, variationImageResp)
	if err != nil {
		return "", err
	}
	var imageBase64 = ""
	for _, content := range variationImageResp.ImageContents {
		imageBase64 = content.B64Json
		break
	}
	return imageBase64, nil

}

func buildRequest(pngSrc *os.File) (*http.Request, error) {
	// Create multipart writer
	byteBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(byteBuf)
	//write normal param
	writer.WriteField("n", "1")
	writer.WriteField("size", dto.IMAGE_SIZE_1024)
	writer.WriteField("response_format", dto.IMAGE_FROMAT_BASE64)
	// Add pngSrc to request
	fileWriter, err := writer.CreateFormFile("image", pngSrc.Name())
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fileWriter, pngSrc)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	// Close multipart writer
	writer.Close()
	req, err := http.NewRequest("POST", config.BASEURL+"images/variations", byteBuf)
	if err != nil {
		return nil, err
	}
	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
