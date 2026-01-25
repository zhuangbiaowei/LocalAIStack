package llm

import (
	"context"
	"fmt"
	"time"
)

type EinoConfig struct {
	Model   string
	Timeout time.Duration
}

type EinoProvider struct {
	cfg EinoConfig
}

func NewEinoProvider(cfg EinoConfig) *EinoProvider {
	return &EinoProvider{cfg: cfg}
}

func (p *EinoProvider) Name() string {
	return "eino"
}

func (p *EinoProvider) Generate(ctx context.Context, req Request) (Response, error) {
	_ = ctx
	_ = req
	return Response{}, fmt.Errorf("eino provider not configured")
}
