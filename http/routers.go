package http

import (
	"github.com/gin-gonic/gin"
)

func route(enging *gin.Engine) {
	enging.GET("/", IndexView)
}
