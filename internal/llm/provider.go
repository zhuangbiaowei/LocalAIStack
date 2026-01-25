package llm

import "context"

type Request struct {
	Model   string
	Prompt  string
	Timeout int
}

type Response struct {
	Text string
}

type Provider interface {
	Name() string
	Generate(ctx context.Context, req Request) (Response, error)
}
