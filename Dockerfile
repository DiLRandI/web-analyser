FROM scratch
ARG APP_PORT=80
COPY ./.bin/web-analyser /web-analyser
ENV APP_PORT=${APP_PORT}
ENV GIN_MODE=release
EXPOSE ${APP_PORT}
ENTRYPOINT ["./web-analyser"]