package main

import (
	"github.com/gin-gonic/gin"
)

func (app *AppConfig) router(r *gin.Engine) {
	r.GET("/health", app.health)
	r.POST("/line/callback", app.lineCallback)
	r.POST("/line/webhook", app.lineWebhook)

}
