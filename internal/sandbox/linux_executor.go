package sandbox

import (
	"bytes"
	"context"
	"os/exec"
	"syscall"
	"time"

	"github.com/ton/framework/internal/domain"
)

type LinuxExecutor struct{}

func NewLinuxExecutor() *LinuxExecutor {
	return &LinuxExecutor{}
}

func (e *LinuxExecutor) ExecuteIsolated(ctx context.Context, req domain.SandboxRequest) (*domain.SandboxResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(req.TimeoutSecs)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", req.Command)
	cmd.Dir = req.ProjectPath

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNET,
	}

	if len(req.AllowedEnvs) > 0 {
		cmd.Env = req.AllowedEnvs
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return &domain.SandboxResponse{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration.String(),
	}, nil
}
