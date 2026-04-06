package crd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type manifest struct {
	URL        string
	ServerSide bool
}

var manifests = []manifest{
	{URL: "https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.5.1/standard-install.yaml"},
	{URL: "https://raw.githubusercontent.com/argoproj/argo-rollouts/stable/manifests/crds/rollout-crd.yaml"},
	{URL: "https://raw.githubusercontent.com/argoproj/argo-rollouts/stable/manifests/crds/experiment-crd.yaml"},
	{URL: "https://raw.githubusercontent.com/argoproj/argo-rollouts/stable/manifests/crds/analysis-run-crd.yaml"},
	{URL: "https://raw.githubusercontent.com/argoproj/argo-rollouts/stable/manifests/crds/analysis-template-crd.yaml"},
	{URL: "https://raw.githubusercontent.com/argoproj/argo-rollouts/stable/manifests/crds/cluster-analysis-template-crd.yaml"},
	{URL: "https://github.com/kedacore/keda/releases/download/v2.19.0/keda-2.19.0-crds.yaml", ServerSide: true},
	{URL: "https://raw.githubusercontent.com/external-secrets/external-secrets/v2.2.0/deploy/crds/bundle.yaml", ServerSide: true},
	{URL: "https://github.com/envoyproxy/gateway/releases/download/v1.7.1/envoy-gateway-crds.yaml", ServerSide: true},
	{URL: "https://raw.githubusercontent.com/kubernetes/autoscaler/vpa-release-1.3/vertical-pod-autoscaler/deploy/vpa-v1-crd-gen.yaml", ServerSide: true},
	{URL: "https://github.com/cert-manager/cert-manager/releases/download/v1.20.0/cert-manager.crds.yaml", ServerSide: true},
}

func InstallAll(ctx context.Context) error {
	for _, m := range manifests {
		args := []string{"apply"}
		if m.ServerSide {
			args = append(args, "--server-side")
		}
		args = append(args, "-f", m.URL)

		cmd := exec.CommandContext(ctx, "kubectl", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = io.Discard

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed applying CRD %s: %w", m.URL, err)
		}
	}
	return nil
}
