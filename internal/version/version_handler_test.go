package version

import (
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUnknownVersion(t *testing.T) {
	h := NewVersionHandler()
	if h == nil {
		t.Error("expected handler to be created")
	}

	mux := http.NewServeMux()
	h.Register("/api/v1", mux)

	req, err := http.NewRequest("GET", "/api/v1/version", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", rr.Code)
	}

	if diff := cmp.Diff(rr.Body.String(), "unknown"); diff != "" {
		t.Errorf("unexpected body: %s", diff)
	}
}

func TestVersion(t *testing.T) {
	os.Setenv("MP_BACKEND_VERSION", "1.2.3")

	h := NewVersionHandler()
	if h == nil {
		t.Error("expected handler to be created")
	}

	mux := http.NewServeMux()
	h.Register("/api/v1", mux)

	req, err := http.NewRequest("GET", "/api/v1/version", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", rr.Code)
	}

	if diff := cmp.Diff(rr.Body.String(), "1.2.3"); diff != "" {
		t.Errorf("unexpected body: %s", diff)
	}
}
