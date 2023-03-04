package main

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"os"
)

type AppConfig struct {
	Bot      *linebot.Client
	InfoLog  *log.Logger
	ErrorLog *log.Logger
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
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := AppConfig{
		Bot:      bot,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
	return app
}
