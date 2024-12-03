package blog

import (
	"net/http"
)

type BlogStore interface {
	GetArticles() ([]Article, error)
	GetArticle(articleID string) (*Article, error)
}

type Handler struct {
	store *BlogStore
}

func NewHandler(s *BlogStore) *Handler {
	return &Handler{
		store: s,
	}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc("GET "+prefix+"/blog/article", h.handleGetArticles)
	mux.HandleFunc("GET "+prefix+"/blog/article/{article_id}", h.handleGetArticle)
}

func (h *Handler) handleGetArticles(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetArticle(w http.ResponseWriter, r *http.Request) {

}
