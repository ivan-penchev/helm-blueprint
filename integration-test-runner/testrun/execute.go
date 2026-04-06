package testrun

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"integration-test-runner/cli"
)

type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("test script failed with exit code %d: %v", e.Code, e.Err)
}

func Execute(ctx context.Context, repoRoot string, opts cli.Options) error {
	scriptPath := filepath.Join(repoRoot, "ci", "kind", "test.sh")

	scriptBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", scriptPath, err)
	}

	normalizedScript := bytes.ReplaceAll(scriptBytes, []byte("\r\n"), []byte("\n"))

	tmpScript, err := os.CreateTemp("", "integration-test-runner-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temp script: %w", err)
	}
	defer os.Remove(tmpScript.Name())

	if _, err := tmpScript.Write(normalizedScript); err != nil {
		tmpScript.Close()
		return fmt.Errorf("failed to write temp script: %w", err)
	}

	if err := tmpScript.Close(); err != nil {
		return fmt.Errorf("failed to close temp script: %w", err)
	}

	if err := os.Chmod(tmpScript.Name(), 0o755); err != nil {
		return fmt.Errorf("failed to chmod temp script: %w", err)
	}

	args := make([]string, 0, 4+len(opts.ExtraArgs))
	if opts.SkipTier1 {
		args = append(args, "--skip-tier1")
	}
	if opts.SkipTier2 {
		args = append(args, "--skip-tier2")
	}
	args = append(args, opts.ExtraArgs...)

	cmdArgs := append([]string{tmpScript.Name()}, args...)
	cmd := exec.CommandContext(ctx, "bash", cmdArgs...)
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return &ExitError{Code: ee.ExitCode(), Err: err}
		}
		return fmt.Errorf("failed to execute %s: %w", scriptPath, err)
	}

	return nil
}
