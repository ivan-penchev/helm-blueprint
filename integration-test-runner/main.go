package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"integration-test-runner/cli"
	"integration-test-runner/cluster"
	"integration-test-runner/crd"
	"integration-test-runner/logging"
	"integration-test-runner/testrun"
	"integration-test-runner/workspace"
)

func main() {
	os.Exit(run())
}

func run() int {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		logging.Error(err.Error())
		return 2
	}

	if err := cluster.CheckPrereqs("kind", "kubectl", "helm"); err != nil {
		logging.Error(err.Error())
		return 1
	}

	repoRoot, err := workspace.FindRepoRoot()
	if err != nil {
		logging.Error(err.Error())
		return 1
	}

	if !opts.KeepCluster {
		defer func() {
			logging.Step(fmt.Sprintf("Deleting kind cluster %q...", opts.ClusterName))
			_ = cluster.Delete(opts.ClusterName)
		}()
	} else {
		defer logging.Info(fmt.Sprintf("Keeping cluster %q (use 'kind delete cluster --name %s' to remove)", opts.ClusterName, opts.ClusterName))
	}

	exists, err := cluster.Exists(opts.ClusterName)
	if err != nil {
		logging.Error(err.Error())
		return 1
	}

	if exists {
		logging.Info(fmt.Sprintf("Cluster %q already exists, reusing it", opts.ClusterName))
	} else {
		logging.Step(fmt.Sprintf("Creating kind cluster %q...", opts.ClusterName))
		start := time.Now()
		if err := cluster.Create(opts.ClusterName, 300*time.Second); err != nil {
			logging.Error(err.Error())
			return 1
		}
		logging.Info(fmt.Sprintf("Cluster created in %ds", int(time.Since(start).Seconds())))
	}

	logging.Step("Installing CRDs...")
	start := time.Now()
	if err := crd.InstallAll(context.Background()); err != nil {
		logging.Error(err.Error())
		return 1
	}
	logging.Info(fmt.Sprintf("CRDs installed in %ds", int(time.Since(start).Seconds())))

	logging.Step("Running integration tests...")
	if err := testrun.Execute(context.Background(), repoRoot, opts); err != nil {
		var ee *testrun.ExitError
		if errors.As(err, &ee) {
			return ee.Code
		}
		logging.Error(err.Error())
		return 1
	}

	return 0
}
