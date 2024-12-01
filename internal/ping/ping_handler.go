package ping

import "net/http"

type Handler struct {
}

func NewPingHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/ping", h.handlePing)
}

func (h *Handler) handlePing(r http.ResponseWriter, req *http.Request) {
	r.WriteHeader(http.StatusOK)
	r.Write([]byte("pong"))

}
