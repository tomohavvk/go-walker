package api

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Routes struct {
	logger slog.Logger
}

func NewRoutes(logger slog.Logger) Routes {
	return Routes{logger: logger}
}

func (h Routes) RegisterHTTPRoutes(engine *gin.Engine) {
	engine.GET("/probes/readiness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	engine.GET("/probes/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "live"})
	})

	h.logger.Info("HTTP routes successfully registered")
}
