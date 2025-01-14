package main

import (
	"backend/internal/analytics"
	"backend/internal/blog"
	"backend/internal/blog/blogstore"
	"backend/internal/exam"
	"backend/internal/exam/examstore"
	"backend/internal/github"
	"backend/internal/logger"
	"backend/internal/newsletter"
	"backend/internal/newsletter/newsletterstore"
	"backend/internal/newsletter/newsletterstore/newsletterdatabase"
	"backend/internal/ping"
	"backend/internal/topic"
	"backend/internal/topic/topicstore"
	"backend/internal/version"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	port := mustGetPort()
	mux := http.NewServeMux()

	mDB := ConnectToMongoDB()

	pullBibliothek()

	analyticsHandler := analytics.NewHandler()
	analyticsHandler.Register(apiPrefix, mux)

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

	newsletterDB := newsletterdatabase.NewDB(mDB.Collection("newsletter"))
	newsletterStore := newsletterstore.NewStore(newsletterDB)
	newsletterHandler := newsletter.NewHandler(newsletterStore)
	newsletterHandler.Register(apiPrefix, mux)

	go func() {
		err := http.ListenAndServe(":"+port, corsMiddleware(recoverMiddleware(logger.LogRequest(mux))))
		if err != nil {
			log.Fatalf("error starting server: %v", err)
			return
		}
	}()
	log.Printf("Server started on port %s", port)

	analytics.SendEvent(analytics.Event{
		Name:       "ServerStarted",
		Properties: map[string]interface{}{},
	})

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

func ConnectToMongoDB() *mongo.Database {
	opts := options.Client().ApplyURI(mustGetMongoURL())
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.Background(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	mDB := client.Database(mustGetMongoDBName())

	fmt.Println("Connected to MongoDB")
	return mDB
}

func pullBibliothek() {
	_, err := github.ListFiles(github.BibliothekRepo, "")
	if err != nil {
		err := github.CloneRepo(github.BibliothekRepo, "https://github.com/morgen-pruefung/bibliothek.git")
		if err != nil {
			log.Fatalf("Error cloning repo: %s\n", err)
		}
	} else {
		err := github.PullRepo(github.BibliothekRepo)
		if err != nil {
			log.Fatalf("Error pulling repo: %s\n", err)
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

func mustGetMongoURL() string {
	url := os.Getenv("MONGODB_URI")
	if url == "" {
		log.Fatalf("MONGODB_URI must be set")
	}

	return url
}

func mustGetMongoDBName() string {
	name := os.Getenv("MONGODB_DB_NAME")
	if name == "" {
		log.Fatalf("MONGODB_DB_NAME must be set")
	}

	return name
}

func mustGetPort() string {
	p := os.Getenv("PORT")
	if p == "" {
		return "4242"
	}

	return p
}
