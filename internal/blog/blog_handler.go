package blog

import (
	"encoding/json"
	"log"
	"net/http"
)

type Store interface {
	GetArticles() ([]Article, error)
	GetArticle(articleID string) (*Article, error)
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
	mux.HandleFunc("GET "+prefix+"/blog/article", h.handleGetArticles)
	mux.HandleFunc("GET "+prefix+"/blog/article/{article_id}", h.handleGetArticle)
}

func (h *Handler) handleGetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.store.GetArticles()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error getting articles: %v\n", err)
		return
	}

	data, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error marshalling articles: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) handleGetArticle(w http.ResponseWriter, r *http.Request) {
	articleID := r.PathValue("article_id")
	article, err := h.store.GetArticle(articleID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error getting article: %v\n", err)
		return
	}

	data, err := json.Marshal(article)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error marshalling article: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
