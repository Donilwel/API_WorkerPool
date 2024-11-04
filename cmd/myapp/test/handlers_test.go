package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"workerPool/internal/api"
)

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/ping", nil)
	w := httptest.NewRecorder()

	api.PingHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	body := w.Body.String()
	if body != "PONG" {
		t.Errorf("got %s, want PONG", body)
	}
}
