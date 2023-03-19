package main

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"strings"
)

const (
	prefix             = "小僕人 "
	imageRequestPrefix = "抽圖 "
	helpPrefix         = "--help"
	InvalidType        = 0 // 訊息請求種類:無效
	QuestionType       = 1 // 訊息請求種類:問問題
	ImageType          = 2 // 訊息請求種類:產生圖片
	InfoIntroType      = 3 // 訊息請求種類:功能介紹
	EmptyType          = 4 // 訊息請求種類:空的訊息
)

type MyMessage struct {
	OriginMsg string // 原始訊息
	Type      int    // 訊息請求種類
	Input     string // 請求訊息
	Reply     string
}

func InitMyMessage(msg string, isSingleUser bool) MyMessage {
	myMsg := MyMessage{
		OriginMsg: msg,
		Type:      InvalidType,
	}
	input := myMsg.OriginMsg
	if !isSingleUser {
		if strings.HasPrefix(msg, prefix) {
			input = strings.Split(msg, prefix)[1]
			setMsg(input, &myMsg)
		}
	} else {
		setMsg(input, &myMsg)
	}
	return myMsg
}

func setMsg(input string, myMsg *MyMessage) {
	if input != "" {
		if strings.HasPrefix(input, helpPrefix) {
			// 功能介紹
			// 小僕人 --help
			myMsg.Type = InfoIntroType
		} else if strings.HasPrefix(input, imageRequestPrefix) {
			// 找圖片
			// 小僕人 抽圖 白色貓咪
			myMsg.Type = ImageType
			myMsg.Input = strings.Split(input, imageRequestPrefix)[1]
		} else {
			// 問問題
			// 小僕人 請列出五間餐廳
			myMsg.Type = QuestionType
			myMsg.Input = input
		}
	} else {
		myMsg.Type = EmptyType
	}
}

func (app *AppConfig) GenerateMyMsgWithLineBot(event *linebot.Event) (MyMessage, error) {
	var err error
	var myMsg MyMessage
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		app.InfoLog.Printf(fmt.Sprintf("event.source=%#v, message=%s", event.Source, message.Text))
		myMsg, err = app.GenerateMsg(message.Text, event.Source.Type == linebot.EventSourceTypeUser)
	}
	app.InfoLog.Printf(fmt.Sprintf("myMsg=%#v", myMsg))
	return myMsg, err
}

func (app *AppConfig) GenerateMsg(txt string, isSingleUser bool) (MyMessage, error) {
	var err error
	// 產生訊息
	reply := "你有說話嗎？"
	myMsg := InitMyMessage(txt, isSingleUser)
	switch myMsg.Type {
	case ImageType:
		reply, err = app.OpenAI.GetImage(myMsg.Input)
	case InfoIntroType:
		reply = "1.想產生圖片請輸入：小僕人 抽圖 白色貓咪\n 2.想問問題請輸入：小僕人 請列出五間餐廳"
	case QuestionType:
		reply, err = app.OpenAI.ChatWithChatGPT(myMsg.Input)
	}
	myMsg.Reply = reply
	return myMsg, err
}
