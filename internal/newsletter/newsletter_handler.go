package newsletter

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type Store interface {
	Subscribe(ctx context.Context, e Entry) error
	Unsubscribe(ctx context.Context, e Entry) error
	GetSubscribers(ctx context.Context) ([]Entry, error)
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
	mux.HandleFunc("POST "+prefix+"/newsletter", h.handleSubscribe)
	mux.HandleFunc("DELETE "+prefix+"/newsletter", h.handleUnsubscribe)
}

func (h *Handler) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error reading request body: %v\n", err)
		return
	}

	var entry Entry
	if json.Unmarshal(data, &entry) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error unmarshalling entry: %v\n", err)
		return
	}

	if h.store.Subscribe(r.Context(), entry) != nil {
		if errors.Is(err, ErrEmptyEmail) || errors.Is(err, ErrorInvalidEmail) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error subscribing to newsletter: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error reading request body: %v\n", err)
		return
	}

	var entry Entry
	if json.Unmarshal(data, &entry) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error unmarshalling entry: %v\n", err)
		return
	}

	if h.store.Unsubscribe(r.Context(), entry) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error unsubscribing from newsletter: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
