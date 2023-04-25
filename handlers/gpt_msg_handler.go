package handlers

import (
	"encoding/base64"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/gtp"
	"github.com/869413421/wechatbot/limit"
	"github.com/eatmoreapple/openwechat"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var gptMsgHandlerGroup gptMessageHandlerGroup

func init() {
	gptMsgHandlerGroup = make([]gptMessageHandler, 0, 5)
	gptMsgHandlerGroup = append(gptMsgHandlerGroup, &friendAddGptMessageHandler{})
	gptMsgHandlerGroup = append(gptMsgHandlerGroup, &chatCompleteMessageHandler{})
	gptMsgHandlerGroup = append(gptMsgHandlerGroup, &textEditMessageHandler{})
	gptMsgHandlerGroup = append(gptMsgHandlerGroup, &imageCreateMessageHandler{})
}

func doHandle(msg *openwechat.Message) {
	for _, handler := range gptMsgHandlerGroup {
		support, err := handler.isSupport(msg)
		if err != nil {
			log.Fatalf("gptMessageHandler#isSupport has error %v", err)
			msg.ReplyText("机器人神了，我一会发现了就去修。")
			break
		}
		if support {
			sender, _ := msg.Sender()
			if limit.ShouldLimit(sender.UserName) {
				msg.ReplyText("请求太频繁，请稍后再使")
			} else {
				handler.handle(msg)
			}
			break
		}
	}

}

type gptMessageHandlerGroup []gptMessageHandler
type gptMessageHandler interface {
	isSupport(msg *openwechat.Message) (bool, error)
	handle(msg *openwechat.Message)
}
type friendAddGptMessageHandler struct{}

func (handler *friendAddGptMessageHandler) isSupport(msg *openwechat.Message) (bool, error) {
	return msg.IsFriendAdd(), nil
}
func (handler *friendAddGptMessageHandler) handle(msg *openwechat.Message) {
	if config.LoadConfig().AutoPass {
		_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
		if err != nil {
			log.Fatalf("add friend agree error : %v", err)
			msg.ReplyText("机器人神了，我一会发现了就去修。")
			return
		}
	}

}

type chatCompleteMessageHandler struct{}

func (handler *chatCompleteMessageHandler) isSupport(msg *openwechat.Message) (bool, error) {
	if msg.IsSendByGroup() {
		if msg.IsAt() {
			return isChatCompleteFromMsg(msg, true)
		}
	}
	if msg.IsSendByFriend() {
		return isChatCompleteFromMsg(msg, false)
	}
	return false, nil
}

func (handler *chatCompleteMessageHandler) handle(msg *openwechat.Message) {
	isGroup := msg.IsSendByGroup()
	reqContent, err := buildRequestPurgeContent(msg, isGroup)
	if err != nil {
		log.Fatalf("buildRequestPurgeContent error : %v", err)
	}
	sender, _ := msg.Sender()
	reply, err := gtp.ChatCompletions(reqContent, sender.UserName, sender.EncryChatRoomId, isGroup)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return
	}
	replyText(msg, isGroup, reply)
}

func replyText(msg *openwechat.Message, isGroup bool, reply string) bool {
	if isGroup {
		// 获取@我的用户
		groupSender, err := msg.SenderInGroup()
		if err != nil {
			log.Printf("get sender in group error :%v \n", err)
			return true
		}
		// 回复@我的用户
		reply = strings.TrimSpace(reply)
		reply = strings.Trim(reply, "\n")
		atText := "@" + groupSender.NickName
		replyText := atText + reply
		_, err = msg.ReplyText(replyText)
		if err != nil {
			log.Printf("response group error: %v \n", err)
		}
	} else {
		// 回复用户
		reply = strings.TrimSpace(reply)
		reply = strings.Trim(reply, "\n")
		_, err := msg.ReplyText(reply)
		if err != nil {
			log.Printf("response user error: %v \n", err)
		}
	}
	return false
}

func isChatCompleteFromMsg(msg *openwechat.Message, isGroup bool) (bool, error) {
	requestText, err := buildRequestPurgeContent(msg, isGroup)
	if err != nil {
		return false, err
	}
	isChatComplete := !strings.HasPrefix(requestText, config.LoadConfig().GenerateImageKeyword) &&
		!strings.HasPrefix(requestText, config.LoadConfig().TextEditKeyword)
	return isChatComplete, nil
}

func buildRequestPurgeContent(msg *openwechat.Message, isGroup bool) (string, error) {
	if isGroup {
		sender, err := msg.Sender()
		if err != nil {
			return "", err
		}
		// 替换掉@文本，然后向GPT发起请求
		replaceText := "@" + sender.Self.NickName
		requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
		return requestText, nil
	} else {
		requestText := strings.TrimSpace(msg.Content)
		requestText = strings.Trim(msg.Content, "\n")
		return requestText, nil
	}
}

type textEditMessageHandler struct{}

func (handler *textEditMessageHandler) isSupport(msg *openwechat.Message) (bool, error) {
	if msg.IsSendByGroup() {
		if msg.IsAt() {
			return isTextEditFromMsg(msg, true)
		}
	}
	if msg.IsSendByFriend() {
		// 向GPT发起请求
		return isTextEditFromMsg(msg, false)
	}
	return false, nil
}
func (handler *textEditMessageHandler) handle(msg *openwechat.Message) {
	isGroup := msg.IsSendByGroup()
	reqContent, err := buildRequestPurgeContent(msg, isGroup)
	if err != nil {
		log.Fatalf("buildRequestPurgeContent error : %v", err)
	}
	sender, _ := msg.Sender()
	reply, err := gtp.TextEdit(reqContent, sender.UserName, sender.EncryChatRoomId, isGroup)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return
	}
	replyText(msg, isGroup, reply)
}

func isTextEditFromMsg(msg *openwechat.Message, isGroup bool) (bool, error) {
	requestText, err := buildRequestPurgeContent(msg, isGroup)
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(requestText, config.LoadConfig().TextEditKeyword), nil
}

type imageCreateMessageHandler struct{}

func (handler *imageCreateMessageHandler) isSupport(msg *openwechat.Message) (bool, error) {
	if msg.IsSendByGroup() {
		if msg.IsAt() {
			return isImageCreateFromMsg(msg, true)
		}
	}
	if msg.IsSendByFriend() {
		// 向GPT发起请求
		return isImageCreateFromMsg(msg, false)
	}
	return false, nil
}

func (handler *imageCreateMessageHandler) handle(msg *openwechat.Message) {
	isGroup := msg.IsSendByGroup()
	reqContent, err := buildRequestPurgeContent(msg, isGroup)
	if err != nil {
		log.Fatalf("buildRequestPurgeContent error : %v", err)
	}
	sender, _ := msg.Sender()
	imageBase64, err := gtp.GenerateImage(reqContent, sender.UserName, sender.EncryChatRoomId, isGroup)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return
	}

	imageData, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		log.Printf("base64 decode error: %v \n", err)
		msg.ReplyText("图片解析失败")
		return
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "generate*.png")
	if err != nil {
		log.Printf("create tmp file failed %v \n", err)
		msg.ReplyText("图片解析失败")
		return
	}
	defer os.Remove(tmpFile.Name())
	// 将字节数组写入文件
	//imgData, _, err := image.Decode(bytes.NewReader(imageData))
	err = ioutil.WriteFile(tmpFile.Name(), imageData, 0666)
	if err != nil {
		log.Printf("file  write error: %v \n", err)
		msg.ReplyText("图片解析失败")
		return
	}
	msg.ReplyImage(tmpFile)
}

func isImageCreateFromMsg(msg *openwechat.Message, isGroup bool) (bool, error) {
	requestText, err := buildRequestPurgeContent(msg, isGroup)
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(requestText, config.LoadConfig().GenerateImageKeyword), nil
}
