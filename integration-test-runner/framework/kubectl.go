package framework

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func ResourceExistsByLabelOrName(ctx context.Context, repoRoot string, namespace string, kind string, release string) (bool, error) {
	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", release)
	out, err := RunCommand(ctx, repoRoot, nil, "kubectl", "get", kind, "-n", namespace, "-l", labelSelector, "-o", "name")
	if err == nil && strings.TrimSpace(out.Stdout) != "" {
		return true, nil
	}

	all, errAll := RunCommand(ctx, repoRoot, nil, "kubectl", "get", kind, "-n", namespace, "-o", "name")
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
	_, err := RunCommand(ctx, repoRoot, nil, "kubectl", "wait", "--for=condition=complete", "job", "-n", namespace, "-l", labelSelector, "--timeout", timeout.String())
	return err
}

func PodsReadyOrComplete(ctx context.Context, repoRoot string, namespace string, release string) (bool, error) {
	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", release)
	out, err := RunCommand(ctx, repoRoot, nil, "kubectl", "get", "pods", "-n", namespace, "-l", labelSelector, "--no-headers")
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

	if events, err := RunCommand(ctx, repoRoot, nil, "kubectl", "get", "events", "-n", namespace, "--sort-by=.lastTimestamp"); err == nil {
		chunks = append(chunks, "=== Events ===\n"+events.Stdout)
	}
	if all, err := RunCommand(ctx, repoRoot, nil, "kubectl", "get", "all", "-n", namespace, "-o", "wide"); err == nil {
		chunks = append(chunks, "=== Resources ===\n"+all.Stdout)
	}

	return strings.Join(chunks, "\n\n")
}
