package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Starting server...")

	port := mustGetPort()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

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
