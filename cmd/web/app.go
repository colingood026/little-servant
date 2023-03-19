package main

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"os"
)

type AppConfig struct {
	LineBot     *linebot.Client
	TelegramBot *tgbotapi.BotAPI
	Env         *EnvKey
	OpenAI      *OpenAI
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
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

// TelegramBotConfig tg bot 所需要的設定
type TelegramBotConfig struct {
	Token        string
	WebhookToken string
}

// EnvKey 系統啟動需要的環境變數
type EnvKey struct {
	LineBot     LineBotConfig
	OpenAI      OpenAiConfig
	TelegramBot TelegramBotConfig
}

const (
	KeyChannelSecret           = "CHANNEL_SECRET"
	keyChannelAccessToken      = "CHANNEL_ACCESS_TOKEN"
	KeyOpenApi                 = "OPEN_API_KEY"
	KeyTelegramBotToken        = "TG_BOT_KEY"
	KeyTelegramBotWebhookToken = "TG_BOT_WEBHOOK_TOKEN"
)

func getEnvKey() (EnvKey, error) {
	var env EnvKey
	// line bot
	cs, ok := os.LookupEnv(KeyChannelSecret)
	if !ok || cs == "" {
		return env, errors.New("can not get env:" + KeyChannelSecret)
	}
	cat, ok := os.LookupEnv(keyChannelAccessToken)
	if !ok || cat == "" {
		return env, errors.New("can not get env:" + keyChannelAccessToken)
	}
	env.LineBot = LineBotConfig{ChannelSecret: cs, ChannelAccessToken: cat}
	// open ai
	openApiKey, ok := os.LookupEnv(KeyOpenApi)
	if !ok || openApiKey == "" {
		return env, errors.New("can not get env:" + KeyOpenApi)
	}
	env.OpenAI = OpenAiConfig{ApiKey: openApiKey}
	// telegram bot
	tgBotToken, ok := os.LookupEnv(KeyTelegramBotToken)
	if !ok || tgBotToken == "" {
		return env, errors.New("can not get env:" + KeyTelegramBotToken)
	}
	tgBotWebhookToken, ok := os.LookupEnv(KeyTelegramBotWebhookToken)
	if !ok || tgBotWebhookToken == "" {
		return env, errors.New("can not get env:" + KeyTelegramBotWebhookToken)
	}
	env.TelegramBot = TelegramBotConfig{
		Token:        tgBotToken,
		WebhookToken: tgBotWebhookToken,
	}
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
	lineBot, err := linebot.New(env.LineBot.ChannelSecret, env.LineBot.ChannelAccessToken)
	if err != nil {
		errorLog.Fatal("initial line bot error:", err)
	}
	// open api
	openAI := NewOpenAI(env.OpenAI.ApiKey)
	// telegram bot
	tgBot, err := tgbotapi.NewBotAPI(env.TelegramBot.Token)
	if err != nil {
		errorLog.Fatal("initial telegram bot error:", err)
	}
	tgBot.Debug = false
	// initial app config
	app := AppConfig{
		LineBot:     lineBot,
		TelegramBot: tgBot,
		Env:         &env,
		OpenAI:      openAI,
		InfoLog:     infoLog,
		ErrorLog:    errorLog,
	}
	return app
}
