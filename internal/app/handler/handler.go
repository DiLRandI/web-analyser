package handler

import (
	"net/http"

	"github.com/DiLRandI/web-analyser/internal/dto"
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
	log.Infof("Processing analysis request")
	req := &dto.AnalysesRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Errorf("Invalid request, %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if req.WebUrl == "" {
		log.Error("WebURL is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusAccepted, &dto.AnalysesResponse{Id: 1})
}

func (h *analysisHandler) getAnalysis(c *gin.Context) {
	log.Infof("Retrieving analysed report")
	c.JSON(http.StatusNoContent, nil)
}
