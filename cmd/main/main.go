package main

import (
	"backend/internal/blog"
	"backend/internal/blog/blogstore"
	"backend/internal/exam"
	"backend/internal/exam/examstore"
	"backend/internal/github"
	"backend/internal/logger"
	"backend/internal/ping"
	"backend/internal/topic"
	"backend/internal/topic/topicstore"
	"backend/internal/version"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	apiPrefix = "/api/v1"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Starting server...")

	port := mustGetPort()
	mux := http.NewServeMux()

	pullBibliothek()

	pingHandler := ping.NewPingHandler()
	pingHandler.Register(apiPrefix, mux)

	versionHandler := version.NewVersionHandler()
	versionHandler.Register(apiPrefix, mux)

	blogStore := blogstore.NewStore()
	blogHandler := blog.NewHandler(blogStore)
	blogHandler.Register(apiPrefix, mux)

	topicStore := topicstore.NewStore()
	topicHandler := topic.NewHandler(topicStore)
	topicHandler.Register(apiPrefix, mux)

	examStore := examstore.NewStore()
	examHandler := exam.NewHandler(examStore)
	examHandler.Register(apiPrefix, mux)

	go func() {
		err := http.ListenAndServe(":"+port, corsMiddleware(recoverMiddleware(logger.LogRequest(mux))))
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
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Middleware to add CORS header
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func pullBibliothek() {
	_, err := github.ListFiles(github.BibliothekRepo, "")
	if err != nil {
		err := github.CloneRepo(github.BibliothekRepo, "https://github.com/morgen-pruefung/bibliothek.git")
		if err != nil {
			log.Fatalf("Error cloning repo: %s\n", err)
		}
	}

	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := github.PullRepo(github.BibliothekRepo)
				if err != nil {
					log.Println("Error pulling repo:", err)
					continue
				}
				log.Printf("Pulled repo %s\n", github.BibliothekRepo)
			}
		}
	}()
}

func mustGetPort() string {
	p := os.Getenv("PORT")
	if p == "" {
		return "4242"
	}

	return p
}
