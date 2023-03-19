package main

import (
	"github.com/gin-gonic/gin"
)

func (app *AppConfig) router(r *gin.Engine) {
	r.GET("/health", app.health)
	r.POST("/line/webhook", app.lineCallback)
	r.POST("/telegram/webhook", app.telegramWebhook)

}
