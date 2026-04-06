package e2e

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type CmdResult struct {
	Stdout string
	Stderr string
}

func runCommand(ctx context.Context, dir string, env []string, name string, args ...string) (CmdResult, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := CmdResult{
		Stdout: strings.TrimSpace(stdout.String()),
		Stderr: strings.TrimSpace(stderr.String()),
	}

	if err != nil {
		return result, fmt.Errorf("%s %s failed: %w\nstdout:\n%s\nstderr:\n%s", name, strings.Join(args, " "), err, result.Stdout, result.Stderr)
	}
	return result, nil
}

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

	_, err := runCommand(ctx, repoRoot, nil, "helm", args...)
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

	_, err := runCommand(ctx, repoRoot, nil, "helm", args...)
	if err != nil {
		return fmt.Errorf("helm install (hookOnly) failed for %s (%s): %w", opts.ReleaseName, opts.ValuesFile, err)
	}
	return nil
}

func HelmUninstall(ctx context.Context, repoRoot string, release string, namespace string) error {
	_, err := runCommand(ctx, repoRoot, nil, "helm", "uninstall", release, "-n", namespace, "--wait", "--timeout", "60s")
	return err
}

func ResourceExistsByLabelOrName(ctx context.Context, repoRoot string, namespace string, kind string, release string) (bool, error) {
	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", release)
	out, err := runCommand(ctx, repoRoot, nil, "kubectl", "get", kind, "-n", namespace, "-l", labelSelector, "-o", "name")
	if err == nil && strings.TrimSpace(out.Stdout) != "" {
		return true, nil
	}

	all, errAll := runCommand(ctx, repoRoot, nil, "kubectl", "get", kind, "-n", namespace, "-o", "name")
	if errAll != nil {
		if err != nil {
			return false, err
		}
		return false, errAll
	}

	for _, line := range strings.Split(all.Stdout, "\n") {
		if strings.Contains(strings.ToLower(line), strings.ToLower(release)) {
			return true, nil
		}
	}
	return false, nil
}

func WaitJobsComplete(ctx context.Context, repoRoot string, namespace string, release string, timeout time.Duration) error {
	if timeout == 0 {
		timeout = 120 * time.Second
	}
	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", release)
	_, err := runCommand(ctx, repoRoot, nil, "kubectl", "wait", "--for=condition=complete", "job", "-n", namespace, "-l", labelSelector, "--timeout", timeout.String())
	return err
}

func PodsReadyOrComplete(ctx context.Context, repoRoot string, namespace string, release string) (bool, error) {
	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", release)
	out, err := runCommand(ctx, repoRoot, nil, "kubectl", "get", "pods", "-n", namespace, "-l", labelSelector, "--no-headers")
	if err != nil {
		return false, err
	}

	lines := strings.Split(strings.TrimSpace(out.Stdout), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "" {
		return false, fmt.Errorf("no pods found for release %s", release)
	}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !(strings.Contains(line, "Running") || strings.Contains(line, "Completed")) {
			return false, nil
		}
	}
	return true, nil
}

func DumpNamespaceDebug(ctx context.Context, repoRoot string, namespace string) string {
	var chunks []string

	if events, err := runCommand(ctx, repoRoot, nil, "kubectl", "get", "events", "-n", namespace, "--sort-by=.lastTimestamp"); err == nil {
		chunks = append(chunks, "=== Events ===\n"+events.Stdout)
	}
	if all, err := runCommand(ctx, repoRoot, nil, "kubectl", "get", "all", "-n", namespace, "-o", "wide"); err == nil {
		chunks = append(chunks, "=== Resources ===\n"+all.Stdout)
	}

	return strings.Join(chunks, "\n\n")
}

var nonDNS1123 = regexp.MustCompile(`[^a-z0-9-]`)

func NewNamespaceName(prefix string, testName string) string {
	basePrefix := sanitizeDNS1123(prefix)
	if basePrefix == "" {
		basePrefix = "hb-it"
	}

	testPart := sanitizeDNS1123(testName)
	if testPart == "" {
		testPart = "case"
	}

	suffix := fmt.Sprintf("%x", time.Now().UnixNano())
	if len(suffix) > 8 {
		suffix = suffix[len(suffix)-8:]
	}

	name := fmt.Sprintf("%s-%s-%s", basePrefix, testPart, suffix)
	if len(name) > 63 {
		name = name[:63]
	}
	name = strings.Trim(name, "-")
	if name == "" {
		return "hb-it-case"
	}
	return name
}

func CreateNamespace(ctx context.Context, repoRoot string, namespace string) error {
	_, err := runCommand(ctx, repoRoot, nil, "kubectl", "create", "namespace", namespace)
	return err
}

func DeleteNamespace(ctx context.Context, repoRoot string, namespace string) error {
	_, err := runCommand(ctx, repoRoot, nil, "kubectl", "delete", "namespace", namespace, "--ignore-not-found=true", "--wait=true", "--timeout=120s")
	return err
}

func sanitizeDNS1123(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	v = strings.ReplaceAll(v, "_", "-")
	v = strings.ReplaceAll(v, "/", "-")
	v = strings.ReplaceAll(v, ".", "-")
	v = nonDNS1123.ReplaceAllString(v, "-")
	v = strings.Trim(v, "-")
	for strings.Contains(v, "--") {
		v = strings.ReplaceAll(v, "--", "-")
	}
	return v
}
