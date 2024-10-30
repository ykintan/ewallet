// main.go
package main

import (
	"ewallet/gateaway/config"
	"ewallet/gateaway/router"
	"ewallet/gateaway/service"
	"log"
)

func main() {
	srv := service.NewServer()
	r := router.SetupRouter(srv)

	log.Println("Starting HTTP server on port", config.GetHTTPPort())
	if err := r.Run(config.GetHTTPPort()); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
