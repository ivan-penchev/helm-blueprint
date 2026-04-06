package main

import (
	"os"

	"integration-test-runner/cli"
	"integration-test-runner/cluster"
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

	if err := cluster.CheckPrereqs("kind", "kubectl", "helm", "go"); err != nil {
		logging.Error(err.Error())
		return 1
	}

	repoRoot, err := workspace.FindRepoRoot()
	if err != nil {
		logging.Error(err.Error())
		return 1
	}

	logging.Step("Running integration tests (testify suite)....")
	if err := testrun.Execute(repoRoot, opts); err != nil {
		logging.Error(err.Error())
		if ee, ok := err.(*testrun.ExitError); ok {
			return ee.Code
		}
		return 1
	}

	return 0
}
