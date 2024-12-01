package main

import (
	"backend/internal/ping"
	"backend/internal/version"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Starting server...")

	port := mustGetPort()
	mux := http.NewServeMux()

	pingHandler := ping.NewPingHandler()
	pingHandler.Register(mux)

	versionHandler := version.NewVersionHandler()
	versionHandler.Register(mux)

	go func() {
		err := http.ListenAndServe(":"+port, mux)
		if err != nil {
			log.Fatalf("error starting server: %v", err)
			return
		}
	}()
	log.Printf("Server started on port %s", port)

	c := make(chan struct{}, 1) // Block forever
	<-c
}

func mustGetPort() string {
	p := os.Getenv("PORT")
	if p == "" {
		return "4242"
	}

	return p
}
