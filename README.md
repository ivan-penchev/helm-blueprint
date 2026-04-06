# helm-blueprint

A single, very opinionated, Helm chart for deploying any Kubernetes application workload.
Instead of maintaining separate charts per application, define your entire deployment through values.

## Quick Start

```bash
helm install my-app ./helm-blueprint -f values.yaml
```

```yaml
containers:
  app:
    image:
      repository: nginx
      tag: "1.27"
    ports:
      http:
        port: 80

networking:
  services:
    http:
      ports:
        http:
          port: 80
```

This produces a Deployment with 1 replica, a ClusterIP Service, and a ServiceAccount.

## Key Features

**Any workload type** -- Deployment, StatefulSet, DaemonSet, CronJob, Job, or Argo Rollout -- all from a single chart with `workloadType`.

**Full networking stack** -- Services, Ingresses, Gateway API routes, Envoy Gateway policies, NetworkPolicies, and cert-manager Certificates.

**Event-driven autoscaling** -- HPA v2, VPA, and KEDA (ScaledObject + ScaledJob) with any trigger.

**Secrets management** -- ConfigMaps, Secrets with push/pull via ESO, and auto-generated passwords (ArgoCD-safe). Stakater Reloader annotations auto-generated per resource.

**Lifecycle hooks** -- Pre-install/pre-upgrade Jobs for DB migrations, schema setup, etc. Auto-generates both Helm and Argo CD hook annotations. Hooks share the main workload's pod settings (SA, volumes, secrets).

**Argo CD native** -- Auto sync waves for ordered deployment, sync options, and dual hook annotations. Works with Argo CD, Flux, and Helm CLI without configuration changes.

**Schema validation** -- Catches misconfigurations at install time: missing mounts, port mismatches, conflicting settings.

## Documentation

### Core

| Guide | Description |
|-------|-------------|
| [Workloads](docs/workloads.md) | Deployment, StatefulSet, DaemonSet, CronJob, Job, Argo Rollout |
| [Containers](docs/containers.md) | Container spec, env, mounts, health checks, init/sidecar containers |
| [Configuration](docs/configuration.md) | ConfigMaps, Secrets (push/pull via ESO), auto-generated secrets |
| [Persistence](docs/persistence.md) | PVCs, StatefulSet volume claim templates, static PV binding |
| [RBAC](docs/rbac.md) | ServiceAccount, Roles, ClusterRoles, Bindings |
| [Scheduling](docs/scheduling.md) | Node affinity, tolerations, topology spread, priority classes |

### Networking

| Guide | Description |
|-------|-------------|
| [Services & Ingresses](docs/networking.md) | ClusterIP, NodePort, LoadBalancer, headless services, Ingress |
| [Gateway API](docs/routes.md) | HTTPRoute, GRPCRoute, TLSRoute, TCPRoute, UDPRoute, Envoy policies |
| [Certificates](docs/certificates.md) | cert-manager TLS certificates with auto-named secrets |
| [Network Policies](docs/network-policies.md) | Ingress/egress rules, deny-all, namespace/pod/IP selectors |

### Autoscaling & Availability

| Guide | Description |
|-------|-------------|
| [Autoscaling](docs/autoscaling.md) | HPA v2, VPA, KEDA (ScaledObject, ScaledJob), Pod Disruption Budgets |

### GitOps & Deployment

| Guide | Description |
|-------|-------------|
| [Argo CD](docs/argocd.md) | Sync waves, sync options, hook annotations |
| [Hooks](docs/hooks.md) | Pre-install/pre-upgrade Jobs, Argo CD + Flux compatible |

### Advanced

| Guide | Description |
|-------|-------------|
| [Resource Quotas](docs/resource-quotas.md) | LimitRange, ResourceQuota |
| [Extra Resources](docs/advanced.md) | Escape hatch for arbitrary resources, global settings, pod settings |

## Requirements

- Kubernetes >= 1.28
- Helm >= 3.x

The core chart (Deployment, Service, Ingress, ConfigMap, Secret, HPA, PDB, RBAC) has **zero external dependencies**. Optional features require their respective operators:

| Feature | Operator | Values key | Version |
|---------|----------|------------|---------|
| Argo Rollouts | [Argo Rollouts](https://argoproj.github.io/rollouts/) | `workloadType: Rollout` | v1.6+ |
| VPA | [Vertical Pod Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler) | `autoscaling.vpa.enabled` | v1.0+ |
| KEDA autoscaling | [KEDA](https://keda.sh/) | `autoscaling.keda.enabled` | v2.12+ |
| Gateway API | [Gateway API CRDs](https://gateway-api.sigs.k8s.io/) | `networking.gatewayApi.routes` | v1.2+ |
| Envoy policies | [Envoy Gateway](https://gateway.envoyproxy.io/) | `networking.gatewayApi.routes.*.policies.envoy` | v1.0+ |
| TLS certificates | [cert-manager](https://cert-manager.io/) | `config.certificates` | v1.12+ |
| External Secrets | [ESO](https://external-secrets.io/) | `config.secrets.*.pull` / `config.secrets.*.push` / `config.secrets.*.generate` | v0.9+ |

## Examples

Tested configurations in [`ci/`](ci/):

| File | Scenario |
|------|----------|
| [`minimal-values.yaml`](ci/minimal-values.yaml) | Simplest deployment |
| [`deployment-values.yaml`](ci/deployment-values.yaml) | Deployment with ingress and HPA |
| [`statefulset-values.yaml`](ci/statefulset-values.yaml) | StatefulSet with persistence |
| [`daemonset-values.yaml`](ci/daemonset-values.yaml) | DaemonSet with service exposure |
| [`cronjob-values.yaml`](ci/cronjob-values.yaml) | Scheduled batch job |
| [`job-values.yaml`](ci/job-values.yaml) | One-shot job |
| [`rollout-values.yaml`](ci/rollout-values.yaml) | Argo Rollout with canary strategy |
| [`keda-values.yaml`](ci/keda-values.yaml) | KEDA ScaledObject |
| [`scaledjob-values.yaml`](ci/scaledjob-values.yaml) | KEDA ScaledJob |
| [`full-values.yaml`](ci/full-values.yaml) | Every feature exercised |

## Values Reference

See the fully commented [`values.yaml`](values.yaml) for all available options.

## License

MIT License -- see [LICENSE](LICENSE).
