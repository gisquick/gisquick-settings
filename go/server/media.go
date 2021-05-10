package server

import (
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gislab-npo/gisquick-settings/fs"
	"github.com/go-chi/chi"
)

func (s *Server) handleMediaFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		filename := chi.URLParam(r, "*")
		path := filepath.Join(s.config.ProjectsRoot, username, directory, "media", filename)
		http.ServeFile(w, r, path)
	}
}

func (s *Server) handleMediaFileUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")

		mediaDir := filepath.Join(s.config.ProjectsRoot, username, directory, "media")

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

		r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024) // max. 10 MB
		reader := multipart.NewReader(r.Body, boundary)
		// res := make([]string, 1)
		// var res []string = []string{}
		var res []string
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Media upload file error: %s\n", err)
				http.Error(w, "Upload error", http.StatusBadRequest)
				return
			}

			destPath := filepath.Join(mediaDir, part.FileName())
			if err = fs.SaveToFile(part, destPath); err != nil {
				log.Printf("Media upload file error: %s\n", err)
				http.Error(w, "Upload error", http.StatusBadRequest)
				return
			}
			res = append(res, filepath.Join("media", part.FileName()))
			log.Println(destPath)
		}
		w.Write([]byte(strings.Join(res, ",")))
	}
}
