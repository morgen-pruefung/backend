package ping

import (
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	mux := http.NewServeMux()
	h := NewPingHandler()
	h.Register(mux)

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	if diff := cmp.Diff("pong", rr.Body.String()); diff != "" {
		t.Errorf("unexpected body: %s", diff)
	}
}
