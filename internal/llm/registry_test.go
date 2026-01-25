package llm

import (
	"context"
	"testing"
)

type stubProvider struct {
	name string
}

func (s stubProvider) Name() string {
	return s.name
}

func (s stubProvider) Generate(_ context.Context, _ Request) (Response, error) {
	return Response{Text: "ok"}, nil
}

func TestRegistryRegisterAndList(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(stubProvider{name: "b"}); err != nil {
		t.Fatalf("register provider b: %v", err)
	}
	if err := registry.Register(stubProvider{name: "a"}); err != nil {
		t.Fatalf("register provider a: %v", err)
	}

	providers := registry.Providers()
	if len(providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(providers))
	}
	if providers[0] != "a" || providers[1] != "b" {
		t.Fatalf("expected sorted providers [a b], got %v", providers)
	}
}

func TestRegistryDuplicate(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(stubProvider{name: "dup"}); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	if err := registry.Register(stubProvider{name: "dup"}); err == nil {
		t.Fatalf("expected duplicate register error")
	}
}
