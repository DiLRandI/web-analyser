FROM scratch
COPY ./.bin/web-analyser /web-analyser
EXPOSE 8080
ENTRYPOINT ["./web-analyser"]