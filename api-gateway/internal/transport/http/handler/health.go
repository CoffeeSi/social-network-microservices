package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (g *Gateway) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
