{{/*
Validate workload type.
*/}}
{{- define "helm-blueprint.validation.workloads.type" -}}
{{- $allowedTypes := list "Deployment" "StatefulSet" "CronJob" "Job" "DaemonSet" "Rollout" }}
{{- if not (has .Values.workloadType $allowedTypes) }}
{{- fail (printf "workloadType %q is invalid. Must be one of: %s" .Values.workloadType (join ", " $allowedTypes)) }}
{{- end }}
{{- end }}

{{/*
Validate CronJob requires schedule.
*/}}
{{- define "helm-blueprint.validation.workloads.cronjob" -}}
{{- if eq .Values.workloadType "CronJob" }}
{{- if not .Values.workloads.cronJob.schedule }}
{{- fail "cronJob.schedule: required when workloadType is CronJob" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Run all workload validations.
*/}}
{{- define "helm-blueprint.validation.workloads" -}}
{{- include "helm-blueprint.validation.workloads.type" . }}
{{- include "helm-blueprint.validation.workloads.cronjob" . }}
{{- end }}
