FROM golang:1.15-alpine as go-builder

RUN apk add --no-cache git
ENV GOPATH=
COPY ./fs/go.* /go/fs/
COPY ./server/go.* /go/server/
WORKDIR /go/server/
RUN go mod download
COPY ./fs /go/fs
COPY ./server /go/server
RUN go build -ldflags="-s -w" -o /go/bin/server cmd/main.go


FROM alpine:latest

WORKDIR /var/www/gisquick-settings
COPY --from=go-builder /go/bin/server ./server

EXPOSE 8001
CMD ["./server"]
