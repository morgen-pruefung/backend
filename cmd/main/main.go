package main

import (
	"backend/internal/blog"
	"backend/internal/blog/blogstore"
	"backend/internal/logger"
	"backend/internal/ping"
	"backend/internal/version"
	"log"
	"net/http"
	"os"
)

const (
	apiPrefix = "/api/v1"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Starting server...")

	port := mustGetPort()
	mux := http.NewServeMux()

	pingHandler := ping.NewPingHandler()
	pingHandler.Register(apiPrefix, mux)

	versionHandler := version.NewVersionHandler()
	versionHandler.Register(apiPrefix, mux)

	blogStore := blogstore.NewStore()
	blogstore.StartUpdateTicker()
	blogHandler := blog.NewHandler(blogStore)
	blogHandler.Register(apiPrefix, mux)

	go func() {
		err := http.ListenAndServe(":"+port, recoverMiddleware(logger.LogRequest(mux)))
		if err != nil {
			log.Fatalf("error starting server: %v", err)
			return
		}
	}()
	log.Printf("Server started on port %s", port)

	c := make(chan struct{}, 1) // Block forever
	<-c
}

// Middleware to catch panics and return an error response
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Recovered from panic: %v", rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func mustGetPort() string {
	p := os.Getenv("PORT")
	if p == "" {
		return "4242"
	}

	return p
}
