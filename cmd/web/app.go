package main

import (
	"errors"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"os"
)

type AppConfig struct {
	LineBot  *linebot.Client
	OpenAI   *OpenAI
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// LineBotConfig line 機器人所需要的設定
type LineBotConfig struct {
	ChannelSecret      string
	ChannelAccessToken string
}

// OpenAiConfig open ai 所需要的設定
type OpenAiConfig struct {
	ApiKey string
}

// EnvKey 系統啟動需要的環境變數
type EnvKey struct {
	LineBot LineBotConfig
	OpenAI  OpenAiConfig
}

func getEnvKey() (EnvKey, error) {
	var env EnvKey
	cs, ok := os.LookupEnv("channelSecret")
	if !ok || cs == "" {
		return env, errors.New("can not get channelSecret")
	}
	cat, ok := os.LookupEnv("channelAccessToken")
	if !ok || cat == "" {
		return env, errors.New("can not get channelAccessToken")
	}
	env.LineBot = LineBotConfig{ChannelSecret: cs, ChannelAccessToken: cat}
	openApiKey, ok := os.LookupEnv("openApiKey")
	if !ok || openApiKey == "" {
		return env, errors.New("can not get openApiKey")
	}
	env.OpenAI = OpenAiConfig{ApiKey: openApiKey}
	return env, nil
}

func NewAppConfig() AppConfig {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	env, err := getEnvKey()
	if err != nil {
		errorLog.Fatal(err)
	}
	// line bot
	bot, err := linebot.New(env.LineBot.ChannelSecret, env.LineBot.ChannelAccessToken)
	if err != nil {
		errorLog.Fatal("initial line bot error:", err)
	}
	// open api
	openAI := NewOpenAI(env.OpenAI.ApiKey)
	// initial app config
	app := AppConfig{
		LineBot:  bot,
		OpenAI:   openAI,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
	return app
}
