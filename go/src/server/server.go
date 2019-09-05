package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

// Config export
type Config struct {
	ProjectsDirectory string
	Server            string
}

// User export
type User struct {
	Username string `json:"username"`
	IsGuest  bool   `json:"is_guest"`
}

/* Structure for managing websocket connections for concurrent access */
type websocketsMap struct {
	sync.Mutex
	connections map[string]*websocket.Conn
}

func (w *websocketsMap) Set(key string, conn *websocket.Conn) {
	w.Lock()
	w.connections[key] = conn
	w.Unlock()
}

func (w *websocketsMap) Get(key string) *websocket.Conn {
	w.Lock()
	defer w.Unlock()
	return w.connections[key]
}

func newWebsocketsMap() *websocketsMap {
	return &websocketsMap{connections: make(map[string]*websocket.Conn)}
}

// Server export
type Server struct {
	config    Config
	router    *chi.Mux
	upgrader  websocket.Upgrader
	pluginsWs *websocketsMap
	appsWs    *websocketsMap
}

type contextKey string

func (c contextKey) String() string {
	return "gisquick " + string(c)
}

var (
	contextKeyUser = contextKey("user")
)

func (s *Server) jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) sendJSONMessage(ws *websocket.Conn, name string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s:%s", name, jsonData)
	ws.WriteMessage(websocket.TextMessage, []byte(msg))
	return nil
}

func (s *Server) routes() {
	s.router.Get("/ws/plugin/{user}", s.handlePluginWs())
	s.router.Get("/ws/app", s.loginRequired(s.handleAppWs()))
	s.router.Get("/api/project/files/{user}/{directory}", s.loginRequired(s.handleProjectFiles()))
	s.router.Post("/api/project/upload", s.loginRequired(s.handleNewUpload()))
	// TODO: add authentication for upload (plugin)
	s.router.Post("/api/project/upload/{user}/{directory}", s.handleUpload())
	s.router.Get("/api/project/download/{user}/{directory}", s.loginRequired(s.handleDownload()))
	s.router.Delete("/api/project/delete/{user}/{directory}", s.loginRequired(s.handleProjectDelete()))
	s.router.Post("/api/project/config/{user}/{directory}", s.loginRequired(s.handleSaveConfig()))
	s.router.Post("/api/project/meta/{user}/{directory}/{name}", s.loginRequired(s.handleSaveProjectMeta()))
	// s.router.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("web"))))
	s.router.Handle("/static/*", http.FileServer(http.Dir("web")))
	s.router.Handle("/img/*", http.FileServer(http.Dir("web")))
	s.router.Get("/*", s.authMiddleware(s.handleIndex()))
	s.router.Get("/dev/", s.authMiddleware(s.handleDev()))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// NewServer export
func NewServer(config Config) *Server {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	//
	s := Server{config, chi.NewRouter(), upgrader, newWebsocketsMap(), newWebsocketsMap()}
	s.router.Use(middleware.Logger)
	s.routes()
	return &s
}
