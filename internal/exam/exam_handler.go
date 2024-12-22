package exam

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Store interface {
	GetExams() ([]Exam, error)
	GetExam(examID string) (*Exam, error)
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
	mux.HandleFunc("GET "+prefix+"/exams", h.handleGetExams)
	mux.HandleFunc("GET "+prefix+"/exams/{exam_id}", h.handleGetExam)
}

func (h *Handler) handleGetExams(w http.ResponseWriter, r *http.Request) {
	exams, err := h.store.GetExams()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error getting exams: %v\n", err)
		return
	}

	data, err := json.Marshal(exams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error marshalling exams: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(60*60*6))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) handleGetExam(w http.ResponseWriter, r *http.Request) {
	examID := r.PathValue("exam_id")
	exam, err := h.store.GetExam(examID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error getting exam: %v\n", err)
		return
	}

	data, err := json.Marshal(exam)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error marshalling exam: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(60*60*24))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
