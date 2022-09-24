
GO_BUILD_ENVIRONMENT=GOOS=linux CGO_ENABLED=0
GO_BUILD_CMD=$(GO_BUILD_ENVIRONMENT) go build
GO_RUN_CMD=go run
DOCKER_BUILD_CMD=docker build --no-cache
APP_NAME=web-analyser
VERSION=0.1.0
IMAGE_NAME=$(APP_NAME):$(VERSION)
APP_PORT=8080

build:
	$(GO_BUILD_CMD) -o .bin/web-analyser -ldflags="-X 'main.Version=$(VERSION)'" cmd/web-analyser/main.go cmd/web-analyser/version.go

build-image: build
	$(DOCKER_BUILD_CMD) --build-arg APP_PORT=$(APP_PORT) . -t $(IMAGE_NAME)

run:
	$(GO_RUN_CMD) cmd/web-analyser/main.go cmd/web-analyser/version.go

clean:
	rm -rf .bin