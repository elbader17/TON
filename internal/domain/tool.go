package domain

import "context"

type PingRequest struct {
	Message string
}

type PingResponse struct {
	Result string
}

type PingExecutor interface {
	Execute(ctx context.Context, req PingRequest) (*PingResponse, error)
}