package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type analysisHandler struct {
}

func New() *analysisHandler {
	return &analysisHandler{}
}

func (h *analysisHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/analyse", h.analyse)
	router.GET("/analyse", h.getAnalysis)
}

func (h *analysisHandler) analyse(c *gin.Context) {
	log.Infof("Starting analysis for the URL %s", "test")
	c.JSON(http.StatusNoContent, nil)
}

func (h *analysisHandler) getAnalysis(c *gin.Context) {
	log.Infof("Retrieving analysed report")
	c.JSON(http.StatusNoContent, nil)
}
