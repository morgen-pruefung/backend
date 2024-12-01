package ping

import "net/http"

type Handler struct {
}

func NewPingHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/ping", h.handlePing)
}

func (h *Handler) handlePing(r http.ResponseWriter, req *http.Request) {
	r.WriteHeader(http.StatusOK)
	r.Write([]byte("pong"))
}
