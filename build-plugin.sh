#!/bin/sh

rm -rf dist/ && mkdir dist

# Compile binaries locally

# export GOPATH=`pwd`/go/libs:`pwd`/go
# go build -ldflags="-s -w" -buildmode=c-shared -o dist/linux_amd64/gisquick.so go/cmd/plugin/main.go
# CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -buildmode=c-shared -o dist/windows_amd64/gisquick.dll go/cmd/plugin/main.go


# Compile binaries via docker image

export UID=$(id -u)
export GID=$(id -g)
docker run -it --rm \
	-v `pwd`/go:/src/go -v `pwd`/dist/:/dist -v `pwd`/cross-compile.sh:/src/cross-compile.sh \
	--workdir /src \
	--user $UID:$GID \
	dockercore/golang-cross ./cross-compile.sh


# Bundle plugin

cd dist
cp -r ../plugin/ gisquick2
cp linux_amd64/gisquick.so ./gisquick2/
zip -r gisquick2_lin64.zip gisquick2
rm -rf gisquick2

cp -r ../plugin/ gisquick2
cp windows_amd64/gisquick.dll ./gisquick2/
zip -r gisquick2_win64.zip gisquick2
rm -rf gisquick2

cp -r ../plugin/ gisquick2
cp darwin_amd64/gisquick.dylib ./gisquick2/
zip -r gisquick2_darwin64.zip gisquick2
rm -rf gisquick2
