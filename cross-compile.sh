#!/bin/sh

export GOCACHE=/tmp/go-build
export GOPATH=/src/go/libs:/src/go

go get github.com/go-chi/chi
go get github.com/gorilla/websocket

go build -ldflags="-s -w" -buildmode=c-shared -o /dist/linux_amd64/gisquick.so go/cmd/plugin/main.go
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -buildmode=c-shared -o /dist/windows_amd64/gisquick.dll go/cmd/plugin/main.go
CGO_ENABLED=1 CC=o64-clang GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -buildmode=c-shared -o /dist/darwin_amd64/gisquick.dylib go/cmd/plugin/main.go
