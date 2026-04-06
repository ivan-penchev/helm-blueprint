package testcases

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
