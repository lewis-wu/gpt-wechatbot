package gtp

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/dto"
	"github.com/869413421/wechatbot/util"
	"github.com/eatmoreapple/openwechat"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func ImageVariation(msg *openwechat.Message, userName string, groupId string, isGroup bool) (string, error) {
	picture, err := msg.GetPicture()
	if err != nil {
		return "", err
	}
	if picture.StatusCode != 200 {
		return "", errors.New("获取微信图片失败")
	}

	jpgImg, err := jpeg.Decode(picture.Body)
	if err != nil {
		return "", err
	}
	needResize := !(msg.ImgWidth == 1024 && msg.ImgHeight == 1024)

	pngPath, err := util.Jpg2PngAndResize(jpgImg, 1024, 1024, needResize)
	if err != nil {
		return "", err
	}
	pngFile, err := os.Open(pngPath)
	if err != nil {
		return "", err
	}
	defer os.Remove(pngPath)
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
	variationImageResp := &dto.ImageResp{}
	err = json.Unmarshal(body, variationImageResp)
	if err != nil {
		return "", err
	}
	var imageBase64 = ""
	if variationImageResp.Created == 0 {
		log.Printf("GPT ImageVariation error:%v", string(body))
		return "", errors.New("image variation failed")
	}
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
	//req.Header.Set("Content-Length", strconv.Itoa(byteBuf.Len()))
	return req, nil
}
