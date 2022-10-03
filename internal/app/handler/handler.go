package handler

import (
	"net/http"
	"strconv"

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
	router.GET("/analyse/:id", h.getAnalysisById)
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
		return
	}

	c.JSON(http.StatusAccepted, res)
}

func (h *analysisHandler) getAnalysis(c *gin.Context) {
	log.Infof("Retrieving analysed reports")
	res, err := h.processor.GetProcessResults(c)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *analysisHandler) getAnalysisById(c *gin.Context) {
	paramId := c.Param("id")
	if paramId == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		logrus.Errorf("Unable to parse parameter id %q to int, %v", paramId, err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.Infof("Retrieving analysed report for id %d", id)
	res, err := h.processor.GetProcessResultFor(c, id)
	if err != nil {
		if notFoundErr, ok := err.(*service.NotFoundError); ok {
			logrus.Error(notFoundErr)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}
