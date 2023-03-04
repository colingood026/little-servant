package main

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"os"
)

type AppConfig struct {
	Bot *linebot.Client
}

func NewAppConfig() AppConfig {
	cs, ok := os.LookupEnv("channelSecret")
	if !ok || cs == "" {
		log.Fatal("can not get channelSecret")
	}
	cat, ok := os.LookupEnv("channelAccessToken")
	if !ok || cat == "" {
		log.Fatal("can not get channelAccessToken")
	}
	log.Printf("channelSecret %s, channelAccessToken%s, initializing line bot...\n", cs, cat)
	bot, err := linebot.New(cs, cat)
	if err != nil {
		log.Fatal("initial line bot error:", err)
	}
	log.Println("initialized line bot")
	app := AppConfig{
		Bot: bot,
	}
	return app
}
