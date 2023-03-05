package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"net/http"
	"strings"
)

const (
	prefix             = "小僕人 "
	imageRequestPrefix = "抽圖 "
)

func (app *AppConfig) health(c *gin.Context) {
	app.InfoLog.Println("hit health...")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"commit":  Commit,
	})
}

func (app *AppConfig) lineCallback(c *gin.Context) {
	events, err := app.Bot.ParseRequest(c.Request)
	if err != nil {
		// Do something when something bad happened.
		app.ErrorLog.Println(err)
		return
	}
	for _, event := range events {
		app.InfoLog.Printf(fmt.Sprintf("event.source=%v", event.Source))
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				app.InfoLog.Printf(fmt.Sprintf("event.source=%v, message=%s", event.Source, message.Text))
				if strings.HasPrefix(message.Text, prefix) {
					input := strings.Split(message.Text, prefix)[1]
					reply := "你有說話嗎？"
					if input != "" {
						// 找圖片
						// 小僕人 抽圖 白色貓咪
						if strings.HasPrefix(input, imageRequestPrefix) {
							imageDesc := strings.Split(input, imageRequestPrefix)[1]
							reply, err = app.OpenAI.GetImage(imageDesc)
						} else {
							// 問問題
							// 小僕人 請列出五間餐廳
							reply, err = app.OpenAI.ChatWithChatGPT(input)
						}
						if err != nil {
							app.ErrorLog.Println(err)
							reply = "OpenAI 發生錯誤了：" + err.Error()
						}
					}
					if _, err = app.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
						app.ErrorLog.Println(err)
					}
				}
			}
		}
	}
}
