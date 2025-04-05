package webserver

import (
	"log"
	"os"
	"path/filepath"
)

type WebServerConfig struct {
	listenAddr    string
	apiAddr       string
	templatePaths string
}

func DefaultWebServerConfig() *WebServerConfig {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	listenAddr := ":3001"
	apiAddr := "http://localhost:3000/api/v1"
	dir = filepath.Join(dir, "..", "web", "templates", "*.html")
	return &WebServerConfig{
		listenAddr,
		apiAddr,
		dir,
	}
}
