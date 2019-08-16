FROM golang:alpine as builder

RUN apk add --no-cache git

RUN go get github.com/go-chi/chi && \
    go get github.com/gorilla/websocket

COPY src/server /go/src/server
COPY src/fs /go/src/fs
COPY cmd /go/cmd

RUN go build -o /go/bin/server cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /go/bin/server .
COPY web /root/web

EXPOSE 8001

CMD ["./server"]

