package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetMessage (c *gin.Context) {
	message := "Get Message!"
	c.String(http.StatusOK, "%s", message)
}

func PostMessage (c *gin.Context) {
	message := "Posted Message: " + c.PostForm("body")
	c.String(http.StatusOK, "%s", message)
}
