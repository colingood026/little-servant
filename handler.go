package main

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"net/http"
	"strings"
)

const (
	prefix = "小僕人 "
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
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.HasPrefix(message.Text, prefix) {
					input := strings.Split(message.Text, prefix)[1]
					reply := "你有說話嗎？"
					if input != "" {
						reply, err = app.OpenAI.ChatWithChatGPT(input)
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
