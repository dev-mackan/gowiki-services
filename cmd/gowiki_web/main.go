package main

import (
	"github.com/dev-mackan/gowiki/internal/webserver"
	"log"
)

func main() {
	config := webserver.DefaultWebServerConfig()
	web := webserver.NewWebServer(config)
	log.Fatal(web.Run())
}
