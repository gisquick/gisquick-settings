FROM node:12-alpine AS webapp
MAINTAINER Marcel Dancak "dancakm@gmail.com"

WORKDIR /web
COPY ./web/package*.json ./
RUN npm install

COPY ./web ./
RUN npm run build


FROM golang:alpine as goserver

RUN apk add --no-cache git
RUN go get github.com/go-chi/chi && \
    go get github.com/gorilla/websocket

COPY ./go/src /go/src
COPY ./go/cmd /go/cmd
RUN go build -o /go/bin/server cmd/server/main.go


FROM alpine:latest

WORKDIR /root/
COPY --from=goserver /go/bin/server ./server
COPY --from=webapp /web/dist/ ./web

EXPOSE 8001
CMD ["./server"]
