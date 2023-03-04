package main

import (
	"github.com/gin-gonic/gin"
)

func (app *AppConfig) router(r *gin.Engine) {
	r.GET("/health", app.health)
	r.GET("/callback", app.callback)
}
