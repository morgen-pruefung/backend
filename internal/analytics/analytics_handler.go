package analytics

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc("POST "+prefix+"/analytics/page-visited", h.handlePageVisited)
}

func (h *Handler) handlePageVisited(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error reading request body: %v\n", err)
		return
	}

	var pageViewRequest PageVisitRequest
	err = json.Unmarshal(bodyData, &pageViewRequest)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error unmarshalling request body: %v\n", err)
		return
	}

	if pageViewRequest.URL == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	SendEvent(Event{
		Name: "PageVisited",
		Properties: map[string]interface{}{
			"url":     pageViewRequest.URL,
			"referer": pageViewRequest.Referer,
		},
	})

	w.WriteHeader(http.StatusNoContent)
}
