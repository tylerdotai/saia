package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

type Executor struct {
	Shell   string
	CWD     string
	Timeout time.Duration
}

func NewExecutor(shell, cwd string, timeoutSecs int) *Executor {
	return &Executor{
		Shell:   shell,
		CWD:     cwd,
		Timeout: time.Duration(timeoutSecs) * time.Second,
	}
}

func (e *Executor) Run(ctx context.Context, command string) (*Result, error) {
	ctx, cancel := context.WithTimeout(ctx, e.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.Shell, "-c", command)
	cmd.Dir = e.CWD

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
			return nil, fmt.Errorf("exec error: %w", err)
		}
	}

	return &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: duration,
	}, nil
}
