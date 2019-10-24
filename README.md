# gisquick-settings

Server, plugin and web app for publishing QGIS projects in Gisquick

## Build docker image of server
```
docker build -t gisquick/settings .
```


### Build plugin's library

```
export GOPATH=`pwd`/go/libs:`pwd`/go
go build -buildmode=c-shared -o plugin/gisquick.so go/cmd/plugin/main.go
```

#### Windows 64bit (cross compilation on Linux)

Install required packages
```
apt-get install gcc-mingw-w64
```

Then build
```
export GOPATH=`pwd`/go/libs:`pwd`/go
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -o plugin/gisquick.dll go/cmd/plugin/main.go
```

### Plugin development (Linux):
```
ln -s `pwd`/plugin ~/.local/share/QGIS/QGIS3/profiles/default/python/plugins/gisquick2
```
