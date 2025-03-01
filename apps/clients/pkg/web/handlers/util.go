package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/clients/config"
)

type response struct {
	ApiVersion string      `json:"version"`
	Timestamp  time.Time   `json:"timestamp"`
	Status     int         `json:"status"`
	Data       interface{} `json:"data"`
}

func resp(c *gin.Context, code int, obj interface{}) {
	c.JSON(code, response{
		ApiVersion: config.GetVersion(),
		Timestamp:  time.Now(),
		Status:     code,
		Data:       obj,
	})
}
