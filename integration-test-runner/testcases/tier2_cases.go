package testcases

import "time"

type WaitCondition string

const (
	WaitNone     WaitCondition = "none"
	WaitReady    WaitCondition = "ready"
	WaitComplete WaitCondition = "complete"
)

type Tier2Case struct {
	Name         string
	ValuesFile   string
	WaitFor      WaitCondition
	HelmTimeout  time.Duration
	WaitTimeout  time.Duration
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
