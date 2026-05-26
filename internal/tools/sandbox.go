package tools

import (
	"context"

	"github.com/ton/framework/internal/domain"
	"github.com/ton/framework/internal/sandbox"
	"github.com/ton/framework/pkg/errors"
)

type SandboxTool struct {
	executor domain.SandboxExecutor
}

func NewSandboxTool() *SandboxTool {
	return &SandboxTool{executor: sandbox.NewLinuxExecutor()}
}

func (t *SandboxTool) Name() string {
	return "execute_in_sandbox"
}

func (t *SandboxTool) Description() string {
	return "Executes commands in an isolated sandbox environment with network isolation. Use this tool to safely run builds, tests, and other commands that may have side effects. The sandbox provides CLONE_NEWNET network isolation to prevent data exfiltration."
}

func (t *SandboxTool) Execute(ctx context.Context, req interface{}) (interface{}, *errors.TONError) {
	sandboxReq, ok := req.(domain.SandboxRequest)
	if !ok {
		return nil, errors.NewValidation("invalid request type: expected SandboxRequest")
	}
	resp, err := t.executor.ExecuteIsolated(ctx, sandboxReq)
	if err != nil {
		return nil, errors.NewExecution(err.Error())
	}
	return resp, nil
}
