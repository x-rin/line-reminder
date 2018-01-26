package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, status string)  {
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
