package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"net/http"
)

func (app *AppConfig) health(c *gin.Context) {
	app.InfoLog.Println("hit health...")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"commit":  Commit,
	})
}

func (app *AppConfig) lineCallback(c *gin.Context) {
	events, err := app.LineBot.ParseRequest(c.Request)
	if err != nil {
		// Do something when something bad happened.
		app.ErrorLog.Println(err)
		return
	}
	for _, event := range events {
		app.InfoLog.Printf(fmt.Sprintf("event.source=%v", event.Source))
		if event.Type == linebot.EventTypeMessage {
			myMsg, err := app.GenerateMyMsg(event)
			// 準備發送訊息
			if myMsg.Type != InvalidType {
				if err != nil {
					app.ErrorLog.Println(err)
					myMsg.Reply = "OpenAI 發生錯誤了：" + err.Error()
					if _, err = app.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(myMsg.Reply)).Do(); err != nil {
						app.ErrorLog.Println(err)
					}
				} else {
					if myMsg.Type == ImageType {
						if _, err = app.LineBot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(myMsg.Reply, myMsg.Reply)).Do(); err != nil {
							app.ErrorLog.Println(err)
						}
					} else {
						if _, err = app.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(myMsg.Reply)).Do(); err != nil {
							app.ErrorLog.Println(err)
						}
					}
				}
			}
		}
	}
}

func (app *AppConfig) telegramWebhook(c *gin.Context) {
	header := c.Request.Header
	app.InfoLog.Println(fmt.Sprintf("received header:%v", header))
	webhookToken := header.Get("X-Telegram-Bot-Api-Secret-Token")
	if app.Env.TelegramBot.WebhookToken != webhookToken {
		app.ErrorLog.Printf("received webhook secret token %s is not correct\n", webhookToken)
		c.JSON(200, "webhook secret is invalid")
		return
	}
	var m map[string]interface{}
	err := c.Bind(&m)
	if err != nil {
		app.ErrorLog.Println(err)
		c.JSON(200, err)
		return
	}
	app.InfoLog.Println(fmt.Sprintf("%v\n", m))
	c.JSON(200, "success")
}

func (app *AppConfig) GenerateMyMsg(event *linebot.Event) (MyMessage, error) {
	var err error
	var myMsg MyMessage
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		app.InfoLog.Printf(fmt.Sprintf("event.source=%v, message=%s", event.Source, message.Text))
		// 產生訊息
		reply := "你有說話嗎？"
		myMsg = InitMyMessage(message.Text, event.Source.Type == linebot.EventSourceTypeUser)
		switch myMsg.Type {
		case ImageType:
			reply, err = app.OpenAI.GetImage(myMsg.Input)
		case InfoIntroType:
			reply = "1.想產生圖片請輸入：小僕人 抽圖 白色貓咪\n 2.想問問題請輸入：小僕人 請列出五間餐廳"
		case QuestionType:
			reply, err = app.OpenAI.ChatWithChatGPT(myMsg.Input)
		}
		myMsg.Reply = reply
	}
	app.InfoLog.Printf(fmt.Sprintf("myMsg=%v", myMsg))
	return myMsg, err
}
