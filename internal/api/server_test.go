package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zhuangbiaowei/LocalAIStack/internal/config"
)

type providersPayload struct {
	Default   string   `json:"default"`
	Providers []string `json:"providers"`
}

func TestProvidersHandler(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.LLM.Provider = "eino"

	server := NewServer(cfg, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/providers", nil)
	recorder := httptest.NewRecorder()

	server.server.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload providersPayload
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Default != "eino" {
		t.Fatalf("expected default 'eino', got %q", payload.Default)
	}
	if len(payload.Providers) != 1 || payload.Providers[0] != "eino" {
		t.Fatalf("expected providers ['eino'], got %v", payload.Providers)
	}
}
