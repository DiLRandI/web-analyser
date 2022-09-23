package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type handler struct {
}

func New() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/analyse", h.analyse)
	router.GET("/analyse", h.getAnalysis)
}

func (h *handler) analyse(c *gin.Context) {
	log.Infof("Starting analysis for the URL %s", "test")
	c.JSON(http.StatusNoContent, nil)
}

func (h *handler) getAnalysis(c *gin.Context) {
	log.Infof("Retrieving analysed report")
	c.JSON(http.StatusNoContent, nil)
}
