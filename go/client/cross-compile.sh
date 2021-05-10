#!/bin/sh

export GOCACHE=/tmp/go-build
export GOPATH=

go mod download
go build -ldflags="-s -w" -buildmode=c-shared -o /dist/linux_amd64/gisquick.so cmd/main.go
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -buildmode=c-shared -o /dist/windows_amd64/gisquick.dll cmd/main.go
CGO_ENABLED=1 CC=o64-clang GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -buildmode=c-shared -o /dist/darwin_amd64/gisquick.dylib cmd/main.go
