package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DiLRandI/web-analyser/internal/app/handler"
	"github.com/DiLRandI/web-analyser/internal/repository"
	"github.com/DiLRandI/web-analyser/internal/repository/mem"
	"github.com/DiLRandI/web-analyser/internal/service"
	"github.com/DiLRandI/web-analyser/internal/service/webpage"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	appPort := getApplicationPort()
	log.Infof("Starting web-analyser version %s on port %s", Version, appPort)
	router := gin.Default()
	di := initializeDi()
	registerHandlers(router, di)

	if err := router.Run(fmt.Sprintf(":%s", appPort)); err != nil {
		log.Fatalf("Unable to start the server on port %s, %v", appPort, err)
	}
}

func registerHandlers(router *gin.Engine, di *diRegistry) {
	handler.New(di.processor).RegisterRoutes(router)
}

func initializeDi() *diRegistry {
	resultRepo := mem.NewResultInMemory()
	downloader := webpage.NewDownloader(http.DefaultClient)
	analyserFn := func() webpage.Analyser {
		return webpage.NewAnalyser(http.DefaultClient)
	}
	processor := service.NewProcessor(downloader, analyserFn, resultRepo)

	return &diRegistry{
		resultRepo:    resultRepo,
		downloaderSvc: downloader,
		processor:     processor,

		analyserFn: analyserFn,
	}
}

type diRegistry struct {
	resultRepo    repository.Results
	downloaderSvc webpage.Downloader
	processor     service.Processor

	analyserFn func() webpage.Analyser
}

func getApplicationPort() string {
	p := os.Getenv("APP_PORT")
	if p == "" {
		log.Warnf("Application port `APP_PORT` not specified defaulting to 8080")
		p = "8080"
	}

	return p
}
