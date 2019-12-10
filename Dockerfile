FROM node:12-alpine AS js-builder
MAINTAINER Marcel Dancak "dancakm@gmail.com"

WORKDIR /web
COPY ./web/package*.json ./
RUN npm install

COPY ./web ./
RUN npm run build


FROM golang:alpine as go-builder

RUN apk add --no-cache git
RUN go get github.com/go-chi/chi && \
    go get github.com/gorilla/websocket

COPY ./go/src /go/src
COPY ./go/cmd /go/cmd
RUN go build -ldflags="-s -w" -o /go/bin/server cmd/server/main.go


FROM alpine:latest

WORKDIR /var/www/gisquick-settings
COPY --from=go-builder /go/bin/server ./server
COPY --from=js-builder /web/dist/ ./

VOLUME /var/www/gisquick-settings/static/

EXPOSE 8001
CMD ["./server"]
