package server

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"fs"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

func (s *Server) handlePluginWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		fmt.Printf("Plugin WS: %s", username)
		srcConn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.pluginsWs.Set(username, srcConn)

		if appWs := s.appsWs.Get(username); appWs != nil {
			appWs.WriteMessage(websocket.TextMessage, []byte("PluginConnected"))
		}

		for {
			// Read message from source connection
			msgType, msg, err := srcConn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}

			if appWs := s.appsWs.Get(username); appWs != nil {
				// Write message back to browser
				if err = appWs.WriteMessage(msgType, msg); err != nil {
					break // or better reply with error message?
				}
			} else {
				srcConn.WriteMessage(websocket.TextMessage, []byte("AppDisconnected"))
				fmt.Println("Destination not connected!")
			}
		}
		fmt.Println("Source websocket closed!")
		s.pluginsWs.Set(username, nil)
		if appWs := s.appsWs.Get(username); appWs != nil {
			appWs.WriteMessage(websocket.TextMessage, []byte("PluginDisconnected"))
		}
	}
}

func (s *Server) handleAppWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)

		fmt.Printf("Init WS for user: %v\n", user.Username)
		srcConn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.appsWs.Set(user.Username, srcConn)

		for {
			// Read message from source connection
			msgType, msg, err := srcConn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}

			if pluginWs := s.pluginsWs.Get(user.Username); pluginWs != nil {
				// Write message back to browser
				if err = pluginWs.WriteMessage(msgType, msg); err != nil {
					break // or better reply with error message?
				}
			} else {
				srcConn.WriteMessage(websocket.TextMessage, []byte("PluginDisconnected"))
			}
		}
		s.appsWs.Set(user.Username, nil)
	}
}

func (s *Server) handleProjectFiles() http.HandlerFunc {
	projectsDir := s.config.ProjectsDirectory
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")

		files, err := fs.ListDir(filepath.Join(projectsDir, username, directory))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.jsonResponse(w, files)
	}
}

func (s *Server) handleUpload() http.HandlerFunc {
	type fileUploadProgress struct {
		File     string `json:"file"`
		Progress int    `json:"progress"`
	}

	projectsDir := s.config.ProjectsDirectory

	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")

		ctype, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || ctype != "multipart/form-data" {
			http.Error(w, "Invalid content type", 400)
			return
		}
		boundary, ok := params["boundary"]
		if !ok {
			http.Error(w, http.ErrMissingBoundary.Error(), 400)
			return
		}

		reader := multipart.NewReader(r.Body, boundary)
		part, err := reader.NextPart()
		meta, err := ioutil.ReadAll(part)
		fmt.Printf("Payload: %s", meta)
		uploadProgress := make(map[string]int)
		lastNotification := time.Now()
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				if appWs := s.appsWs.Get(username); appWs != nil {
					s.sendJSONMessage(appWs, "UploadProgress", uploadProgress)
				}
				break
			}
			var partReader io.ReadCloser = part
			if strings.HasSuffix(part.FileName(), ".gz") && !strings.HasSuffix(part.FormName(), ".gz") {
				partReader, _ = gzip.NewReader(part)
			}
			pr := &fs.ProgressReader{Reader: partReader, Step: 32 * 1024, Callback: func(p int) {
				uploadProgress[part.FormName()] = p
				now := time.Now()
				if now.Sub(lastNotification).Seconds() > 0.5 {
					if appWs := s.appsWs.Get(username); appWs != nil {
						s.sendJSONMessage(appWs, "UploadProgress", uploadProgress)
					}
					lastNotification = now
					uploadProgress = make(map[string]int)
				}
			}}
			filename := filepath.Join(projectsDir, username, directory, part.FormName())
			err = fs.SaveToFile(pr, filename)
			partReader.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.Write([]byte(""))
	}
}

func (s *Server) handleNewUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)

		ctype, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || ctype != "multipart/form-data" {
			http.Error(w, "Invalid content type", 400)
			return
		}
		boundary, ok := params["boundary"]
		if !ok {
			http.Error(w, http.ErrMissingBoundary.Error(), 400)
			return
		}

		reader := multipart.NewReader(r.Body, boundary)
		part, _ := reader.NextPart()
		if !strings.HasSuffix(part.FileName(), ".zip") {
			http.Error(w, "Expected zip archive", 400)
			return
		}

		tmpfile, err := ioutil.TempFile("/tmp", part.FileName())
		if err != nil {
			http.Error(w, "FileServer error", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpfile.Name())
		_, err = io.Copy(tmpfile, part)
		if err != nil {
			http.Error(w, "FileServer error", http.StatusInternalServerError)
			return
		}
		archiveReader, err := zip.OpenReader(tmpfile.Name())
		if err != nil {
			http.Error(w, "FileServer error", http.StatusInternalServerError)
			return
		}
		for _, f := range archiveReader.File {
			dest := filepath.Join(s.config.ProjectsDirectory, user.Username, f.Name)
			if !f.FileInfo().IsDir() {
				fr, _ := f.Open()
				err = fs.SaveToFile(fr, dest)
				fr.Close()
				if err != nil {
					http.Error(w, "FileServer error", http.StatusInternalServerError)
					return
				}
			}
		}
		w.Write([]byte(""))
	}
}

func (s *Server) handleDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")

		writer := zip.NewWriter(w)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", directory))
		defer writer.Close()

		projectDir := filepath.Join(s.config.ProjectsDirectory, username, directory)
		projectDir, _ = filepath.Abs(projectDir)
		err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				part, err := writer.Create(path[len(projectDir)+1:])
				if err != nil {
					return err
				}
				if err = fs.CopyFile(part, path); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("FileServer error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleProjectDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		if user.Username != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		dest := filepath.Join(s.config.ProjectsDirectory, username, directory)
		err := os.RemoveAll(dest)
		if err != nil {
			http.Error(w, "FileServer Error", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(""))
	}
}

func (s *Server) handleIndex() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("web/index.html"))

	type AppData struct {
		User *User `json:"user"`
	}
	type TemplateData struct {
		AppData AppData
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data := AppData{}
		user, ok := r.Context().Value(contextKeyUser).(*User)
		if ok && !user.IsGuest {
			data.User = user
		}
		tmpl.Execute(w, TemplateData{data})
	}
}

func (s *Server) handleDev() http.HandlerFunc {
	type Data struct {
		User *User `json:"user"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data := Data{}
		user, ok := r.Context().Value(contextKeyUser).(*User)
		if ok && !user.IsGuest {
			data.User = user
		}
		s.jsonResponse(w, data)
	}
}
