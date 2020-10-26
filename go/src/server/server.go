package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

// Config export
type Config struct {
	ProjectsRoot   string
	MapCacheRoot   string
	AppServer      string
	MapServer      string
	MaxFileUpload  int64
	MaxProjectSize int64
}

// User export
type User struct {
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	IsGuest     bool   `json:"is_guest"`
	IsSuperuser bool   `json:"is_superuser"`
}

type message struct {
	Type   string          `json:"type"`
	Status int             `json:"status,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

type genericMessage struct {
	Type   string      `json:"type"`
	Status int         `json:"status,omitempty"`
	Data   interface{} `json:"data"`
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
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to serialize JSON: %s\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func (s *Server) sendJSONMessage(ws *websocket.Conn, name string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ws.WriteJSON(message{Type: name, Data: jsonData})
}

func (s *Server) apiRoutes() {
	s.router.Get("/ws/plugin", s.loginRequired(s.handlePluginWs()))
	s.router.Get("/ws/app", s.loginRequired(s.handleAppWs()))
	s.router.Get("/api/project/files/{user}/{directory}", s.loginRequired(s.handleProjectFiles()))
	s.router.Post("/api/project/upload", s.loginRequired(s.handleArchiveUpload()))
	s.router.Post("/api/project/upload/{user}/{directory}", s.loginRequired(s.handleUpload()))
	s.router.Get("/api/project/download/{user}/{directory}", s.loginRequired(s.handleDownload()))
	s.router.Delete("/api/project/delete/{user}/{directory}", s.loginRequired(s.handleProjectDelete()))
	s.router.Post("/api/project/config/{user}/{directory}/{name}", s.loginRequired(s.handleSaveConfig()))
	s.router.Post("/api/project/meta/{user}/{directory}/{name}", s.loginRequired(s.handleSaveProjectMeta()))
	s.router.Get("/api/project/meta/{user}/{directory}/{name}", s.loginRequired(s.handleGetProjectMeta()))
	s.router.Delete("/api/project/cache/{user}/{directory}/{name}", s.loginRequired(s.handleCacheDelete()))
	s.router.Get("/api/project/map", s.loginRequired(s.handleGetMap()))

	s.router.Get("/api/project/script/{user}/{directory}", s.loginRequired(s.handleScriptsInfo()))
	s.router.Post("/api/project/script/{user}/{directory}", s.loginRequired(s.handleUploadScript()))
	s.router.Delete("/api/project/script/{user}/{directory}/{module}", s.loginRequired(s.handleDeleteScript()))
	s.router.Get("/api/project/static/{user}/{directory}/*", s.handleStaticFile())
}

func (s *Server) devRoutes() {
	s.router.Post("/api/auth/login/", s.handleProxyRequest())
	s.router.HandleFunc("/api/auth/logout/", s.handleProxyRequest())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// NewServer export
func NewServer(config Config, dev bool) *Server {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	//
	s := Server{config, chi.NewRouter(), upgrader, newWebsocketsMap(), newWebsocketsMap()}
	s.router.Use(middleware.Logger)
	s.apiRoutes()
	if dev {
		s.devRoutes()
	}
	return &s
}
