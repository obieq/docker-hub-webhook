package main

import (
	"log"
	"net/http"

	"github.com/obieq/docker-hub-webhook/config"
	"github.com/obieq/docker-hub-webhook/handlers"
)

func main() {
	port := ":4000"
	h := http.NewServeMux()

	// load config file
	// NOTE: once loaded, config.Cfg() can be invoked wherever and whenever
	config.LoadConfig("./config/config.json")
	cfg := config.Cfg()

	// init handlers
	health := new(handlers.HealthHandler)
	dockerHub := new(handlers.DockerHubHandler)

	// mount health handlers
	h.HandleFunc("/docker-hub/health", health.Handle)

	// mount docker hub handlers
	h.HandleFunc("/docker-hub/webhook", dockerHub.Handle)

	if cfg.TLS != nil {
		log.Printf("starting HTTPS (TLS) server on port %s", port)
		if err := http.ListenAndServeTLS(port, cfg.TLS.CertFile, cfg.TLS.KeyFile, h); err != nil {
			log.Fatal("ListenAndServeTLS: ", err)
		}
	} else {
		log.Printf("starting HTTP server on port %s", port)
		err := http.ListenAndServe(port, h)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
