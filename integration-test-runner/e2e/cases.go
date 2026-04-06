package e2e

import "time"

type Tier1Case struct {
	Name          string
	ValuesFile    string
	ExpectedKinds []string
}

var Tier1Cases = []Tier1Case{
	{Name: "t1-deploy", ValuesFile: "ci/deployment-values.yaml", ExpectedKinds: []string{"deployment", "service", "configmap", "ingress", "horizontalpodautoscaler", "poddisruptionbudget", "serviceaccount"}},
	{Name: "t1-sts", ValuesFile: "ci/statefulset-values.yaml", ExpectedKinds: []string{"statefulset", "service", "serviceaccount"}},
	{Name: "t1-ds", ValuesFile: "ci/daemonset-values.yaml", ExpectedKinds: []string{"daemonset", "service", "serviceaccount"}},
	{Name: "t1-cj", ValuesFile: "ci/cronjob-values.yaml", ExpectedKinds: []string{"cronjob", "secret", "serviceaccount"}},
	{Name: "t1-job", ValuesFile: "ci/job-values.yaml", ExpectedKinds: []string{"job", "serviceaccount"}},
	{Name: "t1-rollout", ValuesFile: "ci/rollout-values.yaml", ExpectedKinds: []string{"rollout", "service", "horizontalpodautoscaler", "poddisruptionbudget", "serviceaccount"}},
	{Name: "t1-keda", ValuesFile: "ci/keda-values.yaml", ExpectedKinds: []string{"deployment", "scaledobject", "service", "serviceaccount"}},
	{Name: "t1-sj", ValuesFile: "ci/scaledjob-values.yaml", ExpectedKinds: []string{"scaledjob", "serviceaccount"}},
	{Name: "t1-full", ValuesFile: "ci/full-values.yaml", ExpectedKinds: []string{"deployment", "service", "configmap", "secret", "ingress", "horizontalpodautoscaler", "poddisruptionbudget", "serviceaccount", "role", "rolebinding", "clusterrole", "clusterrolebinding", "persistentvolumeclaim", "httproute", "grpcroute", "tlsroute", "backendtrafficpolicy", "externalsecret", "certificate", "networkpolicy", "verticalpodautoscaler"}},
}

type WaitCondition string

const (
	WaitNone     WaitCondition = "none"
	WaitReady    WaitCondition = "ready"
	WaitComplete WaitCondition = "complete"
)

type Tier2Case struct {
	Name        string
	ValuesFile  string
	WaitFor     WaitCondition
	HelmTimeout time.Duration
	WaitTimeout time.Duration
}

var Tier2Cases = []Tier2Case{
	{Name: "t2-rollout", ValuesFile: "ci/kind/rollout-values.yaml", WaitFor: WaitNone, HelmTimeout: 120 * time.Second, WaitTimeout: 120 * time.Second},
	{Name: "t2-keda", ValuesFile: "ci/kind/keda-values.yaml", WaitFor: WaitNone, HelmTimeout: 120 * time.Second, WaitTimeout: 120 * time.Second},
	{Name: "t2-sj", ValuesFile: "ci/kind/scaledjob-values.yaml", WaitFor: WaitNone, HelmTimeout: 120 * time.Second, WaitTimeout: 120 * time.Second},
	{Name: "t2-deploy", ValuesFile: "ci/kind/deployment-values.yaml", WaitFor: WaitReady, HelmTimeout: 120 * time.Second, WaitTimeout: 120 * time.Second},
	{Name: "t2-sts", ValuesFile: "ci/kind/statefulset-values.yaml", WaitFor: WaitReady, HelmTimeout: 120 * time.Second, WaitTimeout: 120 * time.Second},
	{Name: "t2-ds", ValuesFile: "ci/kind/daemonset-values.yaml", WaitFor: WaitReady, HelmTimeout: 120 * time.Second, WaitTimeout: 120 * time.Second},
	{Name: "t2-cj", ValuesFile: "ci/kind/cronjob-values.yaml", WaitFor: WaitNone, HelmTimeout: 120 * time.Second, WaitTimeout: 120 * time.Second},
	{Name: "t2-job", ValuesFile: "ci/kind/job-values.yaml", WaitFor: WaitComplete, HelmTimeout: 120 * time.Second, WaitTimeout: 120 * time.Second},
	{Name: "t2-full", ValuesFile: "ci/kind/full-values.yaml", WaitFor: WaitReady, HelmTimeout: 180 * time.Second, WaitTimeout: 120 * time.Second},
}
