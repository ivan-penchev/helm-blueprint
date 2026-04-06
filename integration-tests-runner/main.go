package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"integration-test-runner/cli"
	"integration-test-runner/cluster"
	"integration-test-runner/logging"
	"integration-test-runner/workspace"
)

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		logging.Error(err.Error())
		os.Exit(1)
	}

	if err := cluster.CheckPrereqs("kind", "kubectl", "helm", "go"); err != nil {
		logging.Error(err.Error())
		os.Exit(1)
	}

	repoRoot, err := workspace.FindRepoRoot()
	if err != nil {
		logging.Error(err.Error())
		os.Exit(1)
	}

	logging.Step("Running integration tests (testify suite)....")
	if err := executeE2ETestsFolder(repoRoot, opts); err != nil {
		logging.Error(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func executeE2ETestsFolder(repoRoot string, opts cli.Options) error {
	ctx := context.Background()
	moduleRoot := filepath.Join(repoRoot, "integration-tests-runner")

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
		return err
	}

	return nil
}
