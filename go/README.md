### Setup
```
export GOPATH=`pwd`/libs:`pwd`
go get github.com/go-chi/chi
go get github.com/gorilla/websocket
```

### Run server
```
export GOPATH=`pwd`/libs:`pwd`
export PROJECTS_DIR=./publish/
export SERVER_URL=http://localhost:8000
...
go run cmd/server/main.go -dev
```

### Build plugin
```
export GOPATH=`pwd`/libs:`pwd`
go build -buildmode=c-shared -o gisquick.so cmd/plugin/main.go
```
