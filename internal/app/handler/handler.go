package handler

import (
	"net/http"

	"github.com/DiLRandI/web-analyser/internal/dto"
	"github.com/DiLRandI/web-analyser/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type analysisHandler struct {
	processor service.Processor
}

func New(processor service.Processor) *analysisHandler {
	return &analysisHandler{
		processor: processor,
	}
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

	res, err := h.processor.ProcessPage(c, req)
	if err != nil {
		logrus.Errorf("error while trying to process the page, %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusAccepted, res)
}

func (h *analysisHandler) getAnalysis(c *gin.Context) {
	log.Infof("Retrieving analysed report")
	c.JSON(http.StatusNoContent, nil)
}
