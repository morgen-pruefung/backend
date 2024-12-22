package topic

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Store interface {
	GetTopics() ([]Topic, error)
	GetTopic(topicID string) (*Topic, error)
	GetTextContent(topicID string) ([]byte, error)
}

type Handler struct {
	store Store
}

func NewHandler(s Store) *Handler {
	return &Handler{
		store: s,
	}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc("GET "+prefix+"/topics", h.handleGetTopics)
	mux.HandleFunc("GET "+prefix+"/topics/{topic_id}", h.handleGetTopic)
	mux.HandleFunc("GET "+prefix+"/topics/{topic_id}/text-content", h.handleGetTextContent)
}

func (h *Handler) handleGetTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := h.store.GetTopics()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error getting topics: %v\n", err)
		return
	}

	data, err := json.Marshal(topics)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error marshalling topics: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(60*60*6))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) handleGetTopic(w http.ResponseWriter, r *http.Request) {
	topicID := r.PathValue("topic_id")
	topic, err := h.store.GetTopic(topicID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error getting topic: %v\n", err)
		return
	}

	data, err := json.Marshal(topic)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error marshalling topic: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(60*60*24))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) handleGetTextContent(w http.ResponseWriter, r *http.Request) {
	topicID := r.PathValue("topic_id")

	textContent, err := h.store.GetTextContent(topicID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error getting text content: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(60*60*24))
	w.WriteHeader(http.StatusOK)
	w.Write(textContent)
}
