# gisquick-settings

## Build docker image
```
docker build -t gisquick/settings .
```


### Build plugin (Linux)

export GOPATH=`pwd`/go/libs:`pwd`/go
go build -buildmode=c-shared -o plugin/gisquick.so go/cmd/plugin/main.go


#### Windows 64bit

Install required packages
```
apt-get install gcc-mingw-w64
```

```
export GOPATH=`pwd`/go/libs:`pwd`/go
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -o plugin/gisquick.dll go/cmd/plugin/main.go
```

### Plugin development (Linux):
ln -s `pwd`/plugin ~/.local/share/QGIS/QGIS3/profiles/default/python/plugins/gisquick2
