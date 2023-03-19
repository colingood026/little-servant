package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
		c.JSON(200, "success")
		return
	}
	for _, event := range events {
		app.InfoLog.Printf(fmt.Sprintf("event.source=%#v", event.Source))
		if event.Type == linebot.EventTypeMessage {
			myMsg, err := app.GenerateMyMsgWithLineBot(event)
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
	c.JSON(200, "success")
}

func (app *AppConfig) telegramWebhook(c *gin.Context) {
	header := c.Request.Header
	app.InfoLog.Println(fmt.Sprintf("received header:%#v", header))
	webhookToken := header.Get("X-Telegram-Bot-Api-Secret-Token")
	if app.Env.TelegramBot.WebhookToken != webhookToken {
		app.ErrorLog.Printf("received webhook secret token %s is not correct\n", webhookToken)
		c.JSON(200, "webhook secret is invalid")
		return
	}
	var up tgbotapi.Update
	err := c.Bind(&up)
	if err != nil {
		app.ErrorLog.Println(err)
		newMsg := tgbotapi.NewMessage(up.Message.Chat.ID, "OpenAI 發生錯誤了："+err.Error())
		if _, err = app.TelegramBot.Send(newMsg); err != nil {
			app.ErrorLog.Println(err)
		}
		c.JSON(200, err)
		return
	}
	app.InfoLog.Println(fmt.Sprintf("request body=%#v\n", up))
	// 產生訊息
	myMsg, err := app.GenerateMsg(up.Message.Text, up.Message.Chat.IsPrivate())
	// 準備發送訊息
	if myMsg.Type != InvalidType {
		var newMsg tgbotapi.Chattable
		if err != nil {
			app.ErrorLog.Println(err)
			myMsg.Reply = "OpenAI 發生錯誤了：" + err.Error()
			newMsg = tgbotapi.NewMessage(up.Message.Chat.ID, myMsg.Reply)
		} else {
			if myMsg.Type == ImageType {
				newMsg = tgbotapi.NewPhotoShare(up.Message.Chat.ID, myMsg.Reply)
			} else {
				newMsg = tgbotapi.NewMessage(up.Message.Chat.ID, myMsg.Reply)
			}
		}
		if _, err = app.TelegramBot.Send(newMsg); err != nil {
			app.ErrorLog.Println(err)
		}
	}
	c.JSON(200, "success")
}
