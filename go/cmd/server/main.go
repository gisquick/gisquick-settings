package main

import (
	"log"
	"net/http"
	"os"
	"server"
	"strconv"
	"strings"
	"syscall"
)

func optEnv(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultValue
}

func parseFileSize(value string) int64 {
	unit := 1
	if strings.HasSuffix(value, "M") {
		unit = 1024 * 1024
	} else if strings.HasSuffix(value, "G") {
		unit = 1024 * 1024 * 1024
	}
	num, err := strconv.Atoi(strings.TrimRight(value, "MGB"))
	if err != nil {
		log.Fatal(err)
	}
	return int64(num * unit)
}

func main() {
	config := server.Config{
		ProjectsDirectory: os.Getenv("PROJECTS_DIR"),
		AppServer:         os.Getenv("SERVER_URL"),
		MapServer:         os.Getenv("MAPSERVER_URL"),
		MaxFileUpload:     parseFileSize(optEnv("MAX_FILE_UPLOAD", "100M")),
		MaxProjectSize:    parseFileSize(optEnv("MAX_PROJECT_SIZE", "200M")),
	}

	s := server.NewServer(config, true)
	syscall.Umask(0)
	log.Fatal(http.ListenAndServe(":8001", s))
}
