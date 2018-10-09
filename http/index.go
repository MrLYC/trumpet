package http

import (
	"github.com/gin-gonic/gin"
)

// IndexView : index of auth server
func IndexView(c *gin.Context) {
	c.String(200, "hello trumpet")
}
