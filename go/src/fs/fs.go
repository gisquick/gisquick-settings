package fs

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// File export
type File struct {
	Path  string    `json:"path"`
	Hash  string    `json:"hash"`
	Size  int64     `json:"size"`
	Mtime time.Time `json:"mtime"`
}

// Checksum export
func Checksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha1.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// ListDir export
func ListDir(root string, checksum bool) (*[]File, error) {
	var files []File = []File{}
	excludeExtRegex := regexp.MustCompile(`(?i).*\.(gpkg-wal|gpkg-shm)$`)

	root, _ = filepath.Abs(root)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hash := ""
			if checksum {
				if hash, err = Checksum(path); err != nil {
					return err
				}
			}
			relPath := path[len(root)+1:]
			if !strings.HasPrefix(relPath, ".gisquick") && !strings.HasSuffix(relPath, "~") && !excludeExtRegex.Match([]byte(relPath)) {
				files = append(files, File{relPath, hash, info.Size(), info.ModTime()})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &files, nil
}

// SaveToFile export
func SaveToFile(src io.Reader, filename string) (err error) {
	err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// more verbose but with better errors propagation
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if _, err := io.Copy(file, src); err != nil {
		return err
	}
	return nil
}

// CopyFile export
func CopyFile(dest io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(dest, file)
	return err
}

// ProgressReader export
type ProgressReader struct {
	Reader   io.Reader
	Callback func(int)
	Step     int
	Progress int
	lastVal  int
}

func (r *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.Progress += n
	if r.Progress-r.lastVal >= r.Step || err == io.EOF {
		r.Callback(r.Progress)
		r.lastVal = r.Progress
	}
	return
}
