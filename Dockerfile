# syntax=docker/dockerfile:1

FROM golang:latest as builder

MAINTAINER ozline

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux

RUN go build -o /wechat

RUN chmod -R 777 wechat

EXPOSE 6666

CMD [ "/wechat" ]