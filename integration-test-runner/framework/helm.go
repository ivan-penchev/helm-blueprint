package framework

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
)

type InstallOptions struct {
	ReleaseName string
	Namespace   string
	ValuesFile  string
	Wait        bool
	Timeout     time.Duration
}

func HelmInstall(ctx context.Context, repoRoot string, opts InstallOptions) error {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	args := []string{
		"install", opts.ReleaseName, ".",
		"-n", opts.Namespace,
		"-f", filepath.Join(repoRoot, opts.ValuesFile),
		"--timeout", timeout.String(),
	}
	if opts.Wait {
		args = append(args, "--wait")
	}

	_, err := RunCommand(ctx, repoRoot, nil, "helm", args...)
	if err != nil {
		return fmt.Errorf("helm install failed for %s (%s): %w", opts.ReleaseName, opts.ValuesFile, err)
	}
	return nil
}

func HelmInstallHookOnly(ctx context.Context, repoRoot string, opts InstallOptions) error {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	args := []string{
		"install", opts.ReleaseName, ".",
		"-n", opts.Namespace,
		"-f", filepath.Join(repoRoot, opts.ValuesFile),
		"--wait=hookOnly",
		"--timeout", timeout.String(),
	}

	_, err := RunCommand(ctx, repoRoot, nil, "helm", args...)
	if err != nil {
		return fmt.Errorf("helm install (hookOnly) failed for %s (%s): %w", opts.ReleaseName, opts.ValuesFile, err)
	}
	return nil
}

func HelmUninstall(ctx context.Context, repoRoot string, release string, namespace string) error {
	_, err := RunCommand(ctx, repoRoot, nil, "helm", "uninstall", release, "-n", namespace, "--wait", "--timeout", "60s")
	return err
}
