package main

import (
	"flag"
	"fmt"
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
		ProjectsRoot:   os.Getenv("PROJECTS_ROOT"),
		MapCacheRoot:   os.Getenv("MAP_CACHE_ROOT"),
		AppServer:      os.Getenv("SERVER_URL"),
		MapServer:      os.Getenv("MAPSERVER_URL"),
		MaxFileUpload:  parseFileSize(optEnv("MAX_FILE_UPLOAD", "100M")),
		MaxProjectSize: parseFileSize(optEnv("MAX_PROJECT_SIZE", "200M")),
	}

	devPtr := flag.Bool("dev", false, "development mode")
	portPtr := flag.Int("port", 8001, "port number")
	flag.Parse()

	s := server.NewServer(config, *devPtr)
	syscall.Umask(0)
	address := fmt.Sprintf(":%d", *portPtr)
	log.Fatal(http.ListenAndServe(address, s))
}
