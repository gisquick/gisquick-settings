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

EXPOSE 8001
CMD ["./server"]
