package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter() (*gin.Engine) {
	return gin.Default()
}

func (r Router) Run() {
	e := r.Engine
	e.Group("/api/v1")
	e.GET("/health", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Connection Successful"})
	})
}

func Run() {
	NewRouter().Run()
}