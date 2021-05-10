# gisquick-settings

Server, plugin and web app for publishing QGIS projects in Gisquick

## Build docker image of server
```
cd go
docker build -t gisquick/settings .
```
Image for development:
```
cd go
docker build -f Dockerfile.dev -t gisquick/settings-dev .
```

## Development

### Build plugin's library

```
cd go/client
go build -buildmode=c-shared -o ../../plugin/gisquick.so cmd/main.go
```

### Plugin development (Linux):
```
ln -s `pwd`/plugin ~/.local/share/QGIS/QGIS3/profiles/default/python/plugins/gisquick2
```
