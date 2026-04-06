package testrun

import (
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
	return fmt.Sprintf("go test failed with exit code %d: %v", e.Code, e.Err)
}

func Execute(repoRoot string, opts cli.Options) error {
	ctx := context.Background()
	moduleRoot := filepath.Join(repoRoot, "integration-test-runner")

	args := []string{"test", "./e2e", "-v", "-count=1"}
	args = append(args, opts.ExtraArgs...)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = moduleRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ITR_CLUSTER_NAME=%s", opts.ClusterName),
		fmt.Sprintf("ITR_NAMESPACE_PREFIX=%s", opts.NamespacePrefix),
		fmt.Sprintf("ITR_MAX_PARALLEL=%d", opts.MaxParallel),
		fmt.Sprintf("ITR_KEEP_CLUSTER=%t", opts.KeepCluster),
		fmt.Sprintf("ITR_REPO_ROOT=%s", repoRoot),
		"CGO_ENABLED=0",
	)

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return &ExitError{Code: ee.ExitCode(), Err: err}
		}
		return fmt.Errorf("failed to execute go test: %w", err)
	}

	return nil
}
