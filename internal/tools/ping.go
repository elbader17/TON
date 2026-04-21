package tools

import (
	"context"

	"github.com/ton/framework/internal/domain"
	"github.com/ton/framework/pkg/errors"
)

type PingExecutor struct{}

func NewPingExecutor() *PingExecutor {
	return &PingExecutor{}
}

func (p *PingExecutor) Execute(ctx context.Context, req domain.PingRequest) (*domain.PingResponse, error) {
	return &domain.PingResponse{
		Result: "pong",
	}, nil
}

type PingTool struct {
	executor domain.PingExecutor
}

func NewPingTool() *PingTool {
	return &PingTool{executor: NewPingExecutor()}
}

func (t *PingTool) Name() string {
	return "ping"
}

func (t *PingTool) Description() string {
	return "Ping tool that returns pong"
}

func (t *PingTool) Execute(ctx context.Context, req interface{}) (interface{}, *errors.TONError) {
	pingReq, ok := req.(domain.PingRequest)
	if !ok {
		return nil, errors.NewValidation("invalid request type")
	}
	resp, err := t.executor.Execute(ctx, pingReq)
	if err != nil {
		return nil, errors.NewExecution(err.Error())
	}
	return resp, nil
}