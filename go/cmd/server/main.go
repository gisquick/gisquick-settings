package main

import (
	"log"
	"net/http"
	"os"
	"server"
	"syscall"
)

func main() {
	config := server.Config{
		ProjectsDirectory: os.Getenv("PROJECTS_DIR"),
		Server:            os.Getenv("SERVER_URL"),
	}

	s := server.NewServer(config)
	syscall.Umask(0)
	log.Fatal(http.ListenAndServe(":8001", s))
}
