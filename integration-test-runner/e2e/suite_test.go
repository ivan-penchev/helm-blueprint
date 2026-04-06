package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"integration-test-runner/cluster"
	"integration-test-runner/crd"
	"integration-test-runner/workspace"

	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	cfg := LoadConfig()
	require.NoError(t, cfg.Validate())

	repoRoot := cfg.RepoRoot
	if strings.TrimSpace(repoRoot) == "" {
		var err error
		repoRoot, err = workspace.FindRepoRoot()
		require.NoError(t, err)
	}

	require.NoError(t, cluster.CheckPrereqs("kind", "kubectl", "helm"))

	exists, err := cluster.Exists(cfg.ClusterName)
	require.NoError(t, err)

	if !exists {
		t.Logf("Creating kind cluster %q", cfg.ClusterName)
		require.NoError(t, cluster.Create(cfg.ClusterName, 300*time.Second))
	} else {
		t.Logf("Reusing existing kind cluster %q", cfg.ClusterName)
	}

	t.Log("Installing CRDs")
	require.NoError(t, crd.InstallAll(context.Background()))

	if !cfg.KeepCluster {
		t.Cleanup(func() {
			t.Logf("Deleting kind cluster %q", cfg.ClusterName)
			if err := cluster.Delete(cfg.ClusterName); err != nil {
				t.Logf("cluster cleanup warning: %v", err)
			}
		})
	} else {
		t.Logf("Keeping kind cluster %q", cfg.ClusterName)
	}

	sem := make(chan struct{}, cfg.MaxParallel)

	t.Run("tier1", func(t *testing.T) {
		for _, tc := range Tier1Cases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				sem <- struct{}{}
				defer func() { <-sem }()

				runTier1Case(t, repoRoot, cfg.NamespacePrefix, tc)
			})
		}
	})

	t.Run("tier2", func(t *testing.T) {
		for _, tc := range Tier2Cases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				sem <- struct{}{}
				defer func() { <-sem }()

				runTier2Case(t, repoRoot, cfg.NamespacePrefix, tc)
			})
		}
	})
}

func runTier1Case(t *testing.T, repoRoot string, namespacePrefix string, tc Tier1Case) {
	ctx := context.Background()
	namespace := NewNamespaceName(namespacePrefix, tc.Name)
	release := tc.Name

	require.NoError(t, CreateNamespace(ctx, repoRoot, namespace))

	t.Cleanup(func() {
		if err := HelmUninstall(ctx, repoRoot, release, namespace); err != nil {
			t.Logf("helm uninstall warning (%s/%s): %v", namespace, release, err)
		}
		cleanupNamespace(t, ctx, repoRoot, namespace)
	})

	err := HelmInstallHookOnly(ctx, repoRoot, InstallOptions{
		ReleaseName: release,
		Namespace:   namespace,
		ValuesFile:  tc.ValuesFile,
		Timeout:     60 * time.Second,
	})
	require.NoError(t, err)

	for _, kind := range tc.ExpectedKinds {
		exists, err := ResourceExistsByLabelOrName(ctx, repoRoot, namespace, kind, release)
		require.NoError(t, err, "lookup failed for kind=%s", kind)
		require.True(t, exists, "expected kind %s not found for %s", kind, tc.Name)
	}
}

func runTier2Case(t *testing.T, repoRoot string, namespacePrefix string, tc Tier2Case) {
	ctx := context.Background()
	namespace := NewNamespaceName(namespacePrefix, tc.Name)
	release := tc.Name

	require.NoError(t, CreateNamespace(ctx, repoRoot, namespace))

	t.Cleanup(func() {
		if err := HelmUninstall(ctx, repoRoot, release, namespace); err != nil {
			t.Logf("helm uninstall warning (%s/%s): %v", namespace, release, err)
		}
		cleanupNamespace(t, ctx, repoRoot, namespace)
	})

	waitForInstall := tc.WaitFor != WaitNone
	err := HelmInstall(ctx, repoRoot, InstallOptions{
		ReleaseName: release,
		Namespace:   namespace,
		ValuesFile:  tc.ValuesFile,
		Wait:        waitForInstall,
		Timeout:     tc.HelmTimeout,
	})
	require.NoError(t, err)

	switch tc.WaitFor {
	case WaitReady:
		ready, err := PodsReadyOrComplete(ctx, repoRoot, namespace, release)
		require.NoError(t, err)
		require.True(t, ready, "some pods are not ready/completed")
	case WaitComplete:
		require.NoError(t, WaitJobsComplete(ctx, repoRoot, namespace, release, tc.WaitTimeout))
	}
}

func cleanupNamespace(t *testing.T, ctx context.Context, repoRoot string, namespace string) {
	if t.Failed() {
		debug := DumpNamespaceDebug(ctx, repoRoot, namespace)
		if strings.TrimSpace(debug) != "" {
			t.Log("Failure diagnostics:\n" + debug)
		}
	}

	if err := DeleteNamespace(ctx, repoRoot, namespace); err != nil {
		t.Logf("namespace cleanup warning (%s): %v", namespace, err)
	}
}
