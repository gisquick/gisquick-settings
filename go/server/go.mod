module github.com/gislab-npo/gisquick-settings/server

go 1.15

require (
	github.com/gislab-npo/gisquick-settings/fs v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/gorilla/websocket v1.4.2
	golang.org/x/net v0.0.0-20210508051633-16afe75a6701 // indirect
)

replace github.com/gislab-npo/gisquick-settings/fs => ../fs
