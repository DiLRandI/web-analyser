package main

import (
	"fmt"
	"os"

	"github.com/DiLRandI/web-analyser/internal/app/handler"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	appPort := getApplicationPort()
	log.Infof("Starting web-analyser version %s on port %s", Version, appPort)
	router := gin.Default()
	registerHandlers(router)

	if err := router.Run(fmt.Sprintf(":%s", appPort)); err != nil {
		log.Fatalf("Unable to start the server on port %s, %v", appPort, err)
	}
}

func registerHandlers(router *gin.Engine) {
	handler.New().RegisterRoutes(router)
}

func getApplicationPort() string {
	p := os.Getenv("APP_PORT")
	if p == "" {
		log.Warnf("Application port `APP_PORT` not specified defaulting to 8080")
		p = "8080"
	}

	return p
}
