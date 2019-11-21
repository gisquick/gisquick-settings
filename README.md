# gisquick-settings

Server, plugin and web app for publishing QGIS projects in Gisquick

## Build docker image of server
```
docker build -t gisquick/settings .
```


## Development

### Build plugin's library

```
export GOPATH=`pwd`/go/libs:`pwd`/go
go build -buildmode=c-shared -o plugin/gisquick.so go/cmd/plugin/main.go
```

### Plugin development (Linux):
```
ln -s `pwd`/plugin ~/.local/share/QGIS/QGIS3/profiles/default/python/plugins/gisquick2
```

### Using with local Gisquick deployment

When using plugin with self signed certificate on localhost, setup with https://localhost url will not work.
You can use development proxy server with InsecureSkipVerify option enabled.

InsecureSkipVerify controls whether a client verifies the server's certificate chain and host name.
If InsecureSkipVerify is true, TLS accepts any certificate presented by the server and any host name in that certificate.
In this mode, TLS is susceptible to man-in-the-middle attacks. This should be used only for testing.

(Go 1.12 or higher is required to handle reverse proxy of websocket connections with standard http library)
```
go run go/cmd/dev/proxy.go [-port 8002]
```

Then you can use url of proxy server in plugin configuration.
