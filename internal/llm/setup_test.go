package llm

import (
	"testing"

	"github.com/zhuangbiaowei/LocalAIStack/internal/config"
)

func TestNewRegistryFromConfig(t *testing.T) {
	cfg := config.LLMConfig{
		Provider:       "eino",
		Model:          "test-model",
		TimeoutSeconds: 5,
	}

	registry, err := NewRegistryFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	providers := registry.Providers()
	if len(providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(providers))
	}
	if providers[0] != "eino" {
		t.Fatalf("expected provider 'eino', got %v", providers)
	}
}
