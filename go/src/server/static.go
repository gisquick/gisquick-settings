package server

import (
	"encoding/json"
	"fs"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

type scriptInfo struct {
	Path       string   `json:"path"`
	Components []string `json:"components"`
}

func saveScriptsInfo(path string, data map[string]scriptInfo) error {
	dest, err := os.Create(path)
	if err != nil {
		log.Printf("Failed to create scripts metadata file: %s\n", err)
		return err
	}
	defer dest.Close()

	encoder := json.NewEncoder(dest)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(data); err != nil {
		log.Printf("Failed to encode scripts metadata: %s\n", err)
		return err
	}
	return nil
}

func (s *Server) handleUploadScript() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")

		if !user.IsSuperuser && user.Username != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		projectDir := filepath.Join(s.config.ProjectsRoot, username, directory)
		if _, err := os.Stat(projectDir); os.IsNotExist(err) {
			http.Error(w, "Project directory not found", 400)
			return
		}

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

		r.Body = http.MaxBytesReader(w, r.Body, 4*1024*1024) // max. 4 MB
		reader := multipart.NewReader(r.Body, boundary)
		defer r.Body.Close()

		part, err := reader.NextPart()
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if part.FormName() != "info" {
			log.Println("Missing 'info' form field")
			http.Error(w, "Invalid request data", 400)
			return
		}
		var info scriptInfo
		err = json.NewDecoder(part).Decode(&info)
		if err != nil {
			log.Println("Failed to parse 'info' field:", err)
			http.Error(w, "Invalid metadata content", 400)
			return
		}

		part, _ = reader.NextPart()
		filename := filepath.Join(projectDir, "static", part.FileName())
		if err = fs.SaveToFile(part, filename); err != nil {
			log.Printf("Failed to save script: %s (%s)\n", filename, err)
			http.Error(w, "Failed to save script", http.StatusInternalServerError)
			return
		}

		// var scripts map[string]scriptInfo
		scripts := make(map[string]scriptInfo)
		scriptsFile := filepath.Join(s.config.ProjectsRoot, username, directory, "static", "scripts.json")
		f, err := os.Open(scriptsFile)
		if err == nil {
			json.NewDecoder(f).Decode(&scripts)
		}
		modName := strings.SplitN(filepath.Base(info.Path), ".", 2)[0]
		entry, ok := scripts[modName]
		if ok && entry.Path != info.Path {
			fullpath := filepath.Join(s.config.ProjectsRoot, username, directory, "static", entry.Path)
			log.Println("delete", entry.Path, fullpath)
			if err = os.Remove(fullpath); err != nil {
				log.Printf("Failed to delete old script file: %s (%s)\n", fullpath, err)
			}
		}
		scripts[modName] = info

		if err = saveScriptsInfo(scriptsFile, scripts); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		s.jsonResponse(w, scripts)
	}
}

func (s *Server) handleDeleteScript() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		module := chi.URLParam(r, "module")

		if !user.IsSuperuser && user.Username != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// var scripts []scriptInfo
		scripts := make(map[string]scriptInfo)
		scriptsFile := filepath.Join(s.config.ProjectsRoot, username, directory, "static", "scripts.json")
		f, err := os.Open(scriptsFile)
		if err == nil {
			json.NewDecoder(f).Decode(&scripts)
		}
		entry, ok := scripts[module]
		if !ok {
			log.Printf("Script module does not exist: %s\n", module)
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		path := filepath.Join(s.config.ProjectsRoot, username, directory, "static", entry.Path)

		if err = os.Remove(path); err != nil {
			log.Printf("Failed to delete script file: %s\n", path)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		delete(scripts, module)
		if err = saveScriptsInfo(scriptsFile, scripts); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		s.jsonResponse(w, scripts)
	}
}

func (s *Server) handleScriptsInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")

		if !user.IsSuperuser && user.Username != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		scriptsFile := filepath.Join(s.config.ProjectsRoot, username, directory, "static", "scripts.json")
		/*
			w.Header().Set("Content-Type", "application/json")
			http.ServeFile(w, r, scriptsFile)
		*/
		// var scripts map[string]scriptInfo
		scripts := make(map[string]scriptInfo)
		f, err := os.Open(scriptsFile)
		if err == nil {
			json.NewDecoder(f).Decode(&scripts)
		}
		// if scripts == nil {
		// 	scripts = make(map[string]scriptInfo)
		// }
		s.jsonResponse(w, scripts)
	}
}

func (s *Server) handleStaticFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		filename := chi.URLParam(r, "*")
		path := filepath.Join(s.config.ProjectsRoot, username, directory, "static", filename)
		http.ServeFile(w, r, path)
	}
}
