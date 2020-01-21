## Development

### Setup
```
export GOPATH=`pwd`/libs:`pwd`
go get github.com/go-chi/chi
go get github.com/gorilla/websocket
```

### Run server
```
export GOPATH=`pwd`/libs:`pwd`
export PROJECTS_ROOT=./publish/
export SERVER_URL=http://localhost:8000
...
go run cmd/server/main.go -dev
```

### Build plugin's shared library
```
export GOPATH=`pwd`/libs:`pwd`
go build -buildmode=c-shared -o gisquick.so cmd/plugin/main.go
```
