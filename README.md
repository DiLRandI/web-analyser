# web-analyser [![CodeQL](https://github.com/DiLRandI/web-analyser/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/DiLRandI/web-analyser/actions/workflows/codeql.yml)[![Go build and test](https://github.com/DiLRandI/web-analyser/actions/workflows/go.yml/badge.svg)](https://github.com/DiLRandI/web-analyser/actions/workflows/go.yml)[![Docker deploy status](https://github.com/DiLRandI/web-analyser/actions/workflows/docker.yaml/badge.svg)](https://hub.docker.com/r/deleema1/web-analyser)[![Coverage Status](https://coveralls.io/repos/github/DiLRandI/web-analyser/badge.svg?branch=main)](https://coveralls.io/github/DiLRandI/web-analyser?branch=main)

## How to run the project

This project configured with make script

- To build the code `make build` [output will be in *.bin/web-analyser*]
- To run the code `make run` [default port is 8080] you can specify `APP_PORT` to run on specific port ex `make run APP_PORT=8090`
- To run tests `make test`
- To build docker image `make build-image` by default `APP_PORT` value is exposed from the container

Docker image is published to [Docker hub](https://hub.docker.com/r/deleema1/web-analyser) through CD. You can find releases [here](https://github.com/DiLRandI/web-analyser/releases)

## Running with Docker

- The published image expose on port **80** by default. you can specify different port using `APP_PORT` environment variable.

```docker
docker run -it -e APP_PORT=8080 deleema1/web-analyser
```

or you can map the expose port when you running the container.

```docker
docker run -it -p 8080:80 deleema1/web-analyser
```

## How to use the project

- When you open the project with vscode it will prompt for instal recommended plugin for project.
- in **api** folded of the project root you can see sample request file [analyses.http](https://github.com/DiLRandI/web-analyser/blob/main/api/analyses.http) written from [http-client plugin for vs code](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).

## Running the [web client](https://github.com/DiLRandI/web-analyser-client)

- [web-analyser-client](https://github.com/DiLRandI/web-analyser-client) is a Angular project.
- this also configured with make script, you can simply do `make run` to run the angular application in development mode.
- by default  [web-analyser-client](https://github.com/DiLRandI/web-analyser-client) assumes **web-analyser** is running on port `8080` if this is not the case you need to configure the [web-analyser-client](https://github.com/DiLRandI/web-analyser-client) project [environment file](https://github.com/DiLRandI/web-analyser-client/blob/main/src/environments/environment.ts) to point correct **web-analyser** port.

## Improvements to be made

- Instead of using the standard HTTP tokenizer package, use a library like  [Colly](https://github.com/gocolly/colly)
- Add more unit tests to cover all scenario.
- Refactor analyser code.
- Improve error handling.
- Instead of baking [logrus](https://github.com/sirupsen/logrus) directly, wrap it with an interface allow more fine grain control over the logs.
- Improve overall logs.
- Improve concurrency using channel.
- Make background process more reliable with recovery and retries.
- web client currently use poling request method to get the latest update. this can be improved to web socket. which allow server to push once the analysis is done.
- at the moment cors is configured accept request from any origin, and headers event that are not used. this need to be improved to allow only known headers and origin.
