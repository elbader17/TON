package domain

import "context"

type SandboxRequest struct {
	ProjectPath string
	Command     string
	TimeoutSecs int
	AllowedEnvs []string
}

type SandboxResponse struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration string
}

type SandboxExecutor interface {
	ExecuteIsolated(ctx context.Context, req SandboxRequest) (*SandboxResponse, error)
}
