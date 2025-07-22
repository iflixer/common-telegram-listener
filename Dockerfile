# build environment
FROM golang:1.22 AS build-env
WORKDIR /server
COPY src/go.mod ./
RUN go mod download
COPY src src
WORKDIR /server/src
RUN CGO_ENABLED=0 GOOS=linux go build -o /server/build/app .

FROM alpine:3.15
WORKDIR /app

COPY --from=build-env /server/build/app /app/telegrambot

#ENV GITHUB-SHA=<GITHUB-SHA>

ENTRYPOINT [ "/app/telegrambot" ]
#ENTRYPOINT [ "ls", "-la", "/app/httpserver" ]
