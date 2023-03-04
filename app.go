package main

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"os"
)

type AppConfig struct {
	Bot      *linebot.Client
	OpenAI   *OpenAI
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewAppConfig() AppConfig {
	// line bot
	cs, ok := os.LookupEnv("channelSecret")
	if !ok || cs == "" {
		log.Fatal("can not get channelSecret")
	}
	cat, ok := os.LookupEnv("channelAccessToken")
	if !ok || cat == "" {
		log.Fatal("can not get channelAccessToken")
	}
	log.Println("initializing line bot...")
	bot, err := linebot.New(cs, cat)
	if err != nil {
		log.Fatal("initial line bot error:", err)
	}
	log.Println("initialized line bot")
	// open api
	openApiKey, ok := os.LookupEnv("openApiKey")
	if !ok || openApiKey == "" {
		log.Fatal("can not get openApiKey")
	}
	openAI := NewOpenAI(openApiKey)
	// initial app config
	app := AppConfig{
		Bot:      bot,
		OpenAI:   openAI,
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return app
}
