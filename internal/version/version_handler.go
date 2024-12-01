package version

import (
	"net/http"
	"os"
)

type Handler struct {
}

func NewVersionHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/version", h.handleVersion)
}

func (h *Handler) handleVersion(r http.ResponseWriter, req *http.Request) {
	r.WriteHeader(http.StatusOK)

	ver := os.Getenv("MP_BACKEND_VERSION")
	if ver == "" {
		ver = "unknown"
	}

	r.Write([]byte(ver))
}
