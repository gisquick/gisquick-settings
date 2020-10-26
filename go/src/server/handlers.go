package server

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"fs"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

func extractQgzFile(srcPath, destPath string) error {
	zipReader, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	var projectFile *zip.File
	for _, f := range zipReader.File {
		if strings.HasSuffix(f.Name, ".qgs") {
			projectFile = f
			break
		}
	}
	if projectFile == nil {
		return errors.New("Invalid .qgz archive")
	}
	freader, err := projectFile.Open()
	if err != nil {
		return err
	}
	defer freader.Close()
	return fs.SaveToFile(freader, destPath)
}

func (s *Server) handlePluginWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := user.Username
		log.Printf("Plugin WS: %s \n", username)
		log.Printf("Client: %s\n", r.Header.Get("User-Agent"))
		srcConn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.pluginsWs.Set(username, srcConn)

		if appWs := s.appsWs.Get(username); appWs != nil {
			info := map[string]string{"client": r.Header.Get("User-Agent")}
			appWs.WriteJSON(genericMessage{Type: "PluginStatus", Status: 200, Data: info})
		}

		for {
			// Read message from source connection
			msgType, msg, err := srcConn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}

			if appWs := s.appsWs.Get(username); appWs != nil {
				// Write message back to browser
				if err = appWs.WriteMessage(msgType, msg); err != nil {
					break // or better reply with error message?
				}
			}
		}
		s.pluginsWs.Set(username, nil)
		if appWs := s.appsWs.Get(username); appWs != nil {
			appWs.WriteJSON(genericMessage{Type: "PluginStatus", Status: 503})
		}
	}
}

func (s *Server) handleAppWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)

		log.Printf("App WS: %s\n", user.Username)
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
				log.Println(err)
				break
			}
			if bytes.Compare(msg, []byte("Ping")) == 0 {
				continue
			}

			if pluginWs := s.pluginsWs.Get(user.Username); pluginWs != nil {
				if err = pluginWs.WriteMessage(msgType, msg); err != nil {
					break // or better reply with error message?
				}
			} else {
				srcConn.WriteJSON(genericMessage{Type: "PluginStatus", Status: 503})
			}
		}
		s.appsWs.Set(user.Username, nil)
	}
}

func (s *Server) handleProjectFiles() http.HandlerFunc {
	projectsDir := s.config.ProjectsRoot
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")

		projectDir := filepath.Join(projectsDir, username, directory)
		files, err := fs.ListDir(projectDir, true)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "Project not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
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
	type uploadInfo struct {
		Files []fs.File `json:"files"`
	}

	projectsDir := s.config.ProjectsRoot

	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		projectDir := filepath.Join(projectsDir, username, directory)

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

		user := r.Context().Value(contextKeyUser).(*User)
		if !user.IsSuperuser {
			r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxProjectSize)
		}
		reader := multipart.NewReader(r.Body, boundary)

		// first part should contain upload info
		var info uploadInfo
		part, err := reader.NextPart()
		if err != nil {
			log.Printf("Invalid upload stream: %s\n", err)
			http.Error(w, "Invalid upload stream", http.StatusBadRequest)
			return
		}
		err = json.NewDecoder(part).Decode(&info)
		if err != nil {
			log.Printf("Failed to decode upload metadata: %s\n", err)
			http.Error(w, "Invalid upload stream", http.StatusBadRequest)
			return
		}

		// Check size limit for regular users
		if !user.IsSuperuser {
			filesSizeMap := make(map[string]int64)
			currentFiles, err := fs.ListDir(projectDir, false)
			if err == nil {
				for _, f := range *currentFiles {
					filesSizeMap[f.Path] = f.Size
				}
			} else if !os.IsNotExist(err) {
				log.Printf("Failed to list project files in %s: %s\n", projectDir, err)
			}
			for _, f := range info.Files {
				filesSizeMap[f.Path] = f.Size
			}
			var projectSize int64
			for _, size := range filesSizeMap {
				projectSize += size
			}
			if projectSize > s.config.MaxProjectSize {
				http.Error(w, "Upload size is over limit", http.StatusBadRequest)
				return
			}
		}

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
			filename := filepath.Join(projectDir, part.FormName())
			err = fs.SaveToFile(pr, filename)
			partReader.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// extract qgz project files
		for _, f := range info.Files {
			if strings.HasSuffix(f.Path, ".qgz") {
				qgzFile := filepath.Join(projectDir, f.Path)
				qgsFile := strings.TrimSuffix(qgzFile, "qgz") + "qgs"
				if err = extractQgzFile(qgzFile, qgsFile); err != nil {
					log.Printf("Failed to extract qgis project file: %s (%s)\n", qgzFile, err)
					continue
				}
			}
		}
		w.Write([]byte(""))
	}
}

func (s *Server) handleArchiveUpload() http.HandlerFunc {
	qgisExtRegex := regexp.MustCompile(`(?i).*\.(qgs|qgz)$`)
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)

		if r.ContentLength > s.config.MaxFileUpload {
			log.Printf("Upload error: file size is over limit (user: %s)\n", user.Username)
			http.Error(w, "File size is over limit", http.StatusExpectationFailed)
			return
		}
		ctype, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || ctype != "multipart/form-data" {
			log.Printf("Upload error: invalid content type (user: %s)\n", user.Username)
			http.Error(w, "Invalid content type", 400)
			return
		}
		boundary, ok := params["boundary"]
		if !ok {
			log.Printf("Upload error: Invalid content type (user: %s)\n", user.Username)
			http.Error(w, http.ErrMissingBoundary.Error(), 400)
			return
		}

		if !user.IsSuperuser {
			r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxFileUpload)
		}
		reader := multipart.NewReader(r.Body, boundary)
		part, _ := reader.NextPart()
		if !strings.HasSuffix(part.FileName(), ".zip") {
			log.Printf("Upload error: not a zip archive (user: %s, file: %s)\n", user.Username, part.FileName())
			http.Error(w, "Expected zip archive", 400)
			return
		}

		tmpfile, err := ioutil.TempFile("/tmp", part.FileName())
		if err != nil {
			log.Printf("Upload error: %s\n", err)
			http.Error(w, "FileServer error", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpfile.Name())
		_, err = io.Copy(tmpfile, part)
		if err != nil {
			log.Printf("Upload error: %s\n", err)
			http.Error(w, "FileServer error", http.StatusInternalServerError)
			return
		}
		archiveReader, err := zip.OpenReader(tmpfile.Name())

		if err != nil {
			log.Printf("Upload error: %s\n", err)
			http.Error(w, "FileServer error", http.StatusInternalServerError)
			return
		}

		// Check archive structure - all files in one root directory, with QGIS project
		invalidArchiveHandler := func(msg string) {
			filenames := make([]string, len(archiveReader.File))
			for i, f := range archiveReader.File {
				filenames[i] = f.Name
			}
			log.Printf("Upload error: %s (user: %s)\n", msg, user.Username)
			log.Printf("Archive files: [%s]\n", strings.Join(filenames, ", "))
			http.Error(w, msg, http.StatusBadRequest)
		}
		rootDir, _ := filepath.Split(archiveReader.File[0].Name)
		if len(rootDir) == 0 {
			invalidArchiveHandler("Invalid project archive - not a single directory")
			return
		}
		hasQgisProject := false
		for _, f := range archiveReader.File {
			if !strings.HasPrefix(f.Name, rootDir) {
				invalidArchiveHandler("Invalid project archive - not a single directory")
				return
			}
			if qgisExtRegex.Match([]byte(f.Name)) {
				hasQgisProject = true
			}
		}
		if !hasQgisProject {
			invalidArchiveHandler("Invalid project archive - missing QGIS project")
			return
		}

		// Extract files
		for _, f := range archiveReader.File {
			dest := filepath.Join(s.config.ProjectsRoot, user.Username, f.Name)
			if !f.FileInfo().IsDir() {
				fr, _ := f.Open()
				err = fs.SaveToFile(fr, dest)
				fr.Close()
				if err != nil {
					log.Printf("Upload error: failed to extract archive. %s (user: %s)\n", err, user.Username)
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

		projectDir := filepath.Join(s.config.ProjectsRoot, username, directory)
		projectDir, _ = filepath.Abs(projectDir)
		files, err := fs.ListDir(projectDir, false)
		if err != nil {
			log.Printf("Project download error: %s\n", err)
			http.Error(w, "FileServer error", http.StatusInternalServerError)
			return
		}
		for _, f := range *files {
			part, err := writer.Create(filepath.Join(directory, f.Path))
			if err != nil {
				http.Error(w, "FileServer error", http.StatusInternalServerError)
				return
			}
			if err = fs.CopyFile(part, filepath.Join(projectDir, f.Path)); err != nil {
				http.Error(w, "FileServer error", http.StatusInternalServerError)
				return
			}
		}
	}
}

func (s *Server) handleProjectDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		if !user.IsSuperuser && user.Username != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		dest := filepath.Join(s.config.ProjectsRoot, username, directory)
		if err := os.RemoveAll(dest); err != nil {
			http.Error(w, "FileServer Error", http.StatusInternalServerError)
			return
		}

		// TODO: delete map cache
		w.Write([]byte(""))
	}
}

func (s *Server) handleSaveConfig() http.HandlerFunc {
	var maxBodySize int64 = 1024 * 1024

	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		projectName := chi.URLParam(r, "name")
		if !user.IsSuperuser && user.Username != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		defer r.Body.Close()
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		dest := filepath.Join(s.config.ProjectsRoot, username, directory, ".gisquick", projectName+".json")
		err = os.MkdirAll(filepath.Dir(dest), os.ModePerm)
		if err != nil {
			http.Error(w, "FileServer Error", http.StatusInternalServerError)
			return
		}
		err = ioutil.WriteFile(dest, data, 0644)
		if err != nil {
			log.Printf("Failed to save config file: %s\n", err.Error())
			http.Error(w, "FileServer Error", http.StatusInternalServerError)
			return
		}
		/*
			// pretty JSON
			var out bytes.Buffer
			err = json.Indent(&out, data, "", "  ")
			ioutil.WriteFile(dest, out.Bytes(), 0644)
		*/

		w.Write([]byte(""))
	}
}

func (s *Server) handleSaveProjectMeta() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		projectName := chi.URLParam(r, "name")
		if !user.IsSuperuser && user.Username != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		dest := filepath.Join(s.config.ProjectsRoot, username, directory, projectName+".meta")
		defer r.Body.Close()

		// TODO: create saveConfigFile(data []byte, dest string) function on server or in fs
		data, _ := ioutil.ReadAll(r.Body)
		var out bytes.Buffer
		if err := json.Indent(&out, data, "", "  "); err != nil {
			log.Printf("Failed to format project metadata: %s\n", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if err := ioutil.WriteFile(dest, out.Bytes(), 0644); err != nil {
			log.Printf("Failed to save project metadata: %s\n", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		// save content as it is
		// fs.SaveToFile(r.Body, dest)
		w.Write([]byte(""))
	}
}

func (s *Server) handleGetProjectMeta() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		filename := chi.URLParam(r, "name")

		regexString := fmt.Sprintf(`%s(_(\d{10}))?\.meta$`, regexp.QuoteMeta(filename))
		regex := regexp.MustCompile(regexString)
		var matchedFilename string
		matchedTimestamp := -1

		root := filepath.Join(s.config.ProjectsRoot, username, directory)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".meta") {
				groups := regex.FindStringSubmatch(info.Name())
				if len(groups) == 3 {
					timestamp, _ := strconv.Atoi(groups[2])
					if timestamp > matchedTimestamp {
						matchedFilename = groups[0]
						matchedTimestamp = timestamp
					}
				}
			}
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if matchedFilename == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		jsonContent, err := ioutil.ReadFile(filepath.Join(root, matchedFilename))
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		var meta map[string]interface{}
		if err = json.Unmarshal(jsonContent, &meta); err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		meta["project"] = filepath.Join(username, directory, strings.TrimSuffix(matchedFilename, filepath.Ext(matchedFilename)))
		s.jsonResponse(w, meta)
	}
}

func (s *Server) handleGetMap() http.HandlerFunc {
	client := &http.Client{}
	mapserverPublishDir := "/publish"
	return func(w http.ResponseWriter, r *http.Request) {
		mapParam := r.URL.Query().Get("MAP")
		username := strings.Split(mapParam, "/")[0]
		user := r.Context().Value(contextKeyUser).(*User)
		if !user.IsSuperuser && user.Username != username {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		req, _ := http.NewRequest(http.MethodGet, s.config.MapServer, nil)
		query := r.URL.Query()
		query.Set("MAP", filepath.Join(mapserverPublishDir, mapParam))
		req.URL.RawQuery = query.Encode()
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Mapserver proxy request failed: %s\n", err)
			http.Error(w, "Error", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		// w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
		io.Copy(w, resp.Body)
	}
}

func (s *Server) handleCacheDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKeyUser).(*User)
		username := chi.URLParam(r, "user")
		directory := chi.URLParam(r, "directory")
		projectName := chi.URLParam(r, "name")
		if !user.IsSuperuser && user.Username != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		projectPath := filepath.Join(username, directory, projectName)
		cacheName := fmt.Sprintf("%x", md5.Sum([]byte(projectPath)))
		cacheDir := filepath.Join(s.config.MapCacheRoot, cacheName)
		tmpDir := filepath.Join(s.config.MapCacheRoot, fmt.Sprintf("obsolete_%s_%d", cacheName, time.Now().Unix()))

		if err := os.Rename(cacheDir, tmpDir); err != nil {
			if !os.IsNotExist(err) {
				log.Printf("Failed to delete map cache of project: %s (%s)\n", projectPath, err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
		} else {
			// delete asynchronosly
			go func() {
				fmt.Println("Removing obsole cache directory:", tmpDir)
				if err := os.RemoveAll(tmpDir); err != nil {
					log.Printf("Failed to delete map cache of project: %s (%s)\n", projectPath, err)
				}
			}()
		}
		w.Write([]byte(""))
	}
}

func (s *Server) handleProxyRequest() http.HandlerFunc {
	appServerURL, _ := url.Parse(s.config.AppServer)
	proxy := httputil.NewSingleHostReverseProxy(appServerURL)
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
