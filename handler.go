package main

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
)

func (app *AppConfig) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"commit":  Commit,
	})
}

func (app *AppConfig) lineCallback(c *gin.Context) {
	events, err := app.Bot.ParseRequest(c.Request)
	if err != nil {
		// Do something when something bad happened.
	}
	for _, event := range events {

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = app.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
