{{/*
Run all validations by delegating to component-specific validators.
Called once from the workload template.
*/}}
{{- define "helm-blueprint.validate" -}}
{{- include "helm-blueprint.validation.workloads" . }}
{{- include "helm-blueprint.validation.containers" . }}
{{- include "helm-blueprint.validation.networking" . }}
{{- include "helm-blueprint.validation.autoscaling" . }}
{{- include "helm-blueprint.validation.availability" . }}
{{- include "helm-blueprint.validation.rbac" . }}
{{- include "helm-blueprint.validation.hooks" . }}
{{- end }}
