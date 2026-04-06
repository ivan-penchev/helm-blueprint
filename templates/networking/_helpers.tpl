{{/*
Return the appropriate headless service name for StatefulSets
*/}}
{{- define "helm-blueprint.networking.headlessServiceName" -}}
{{- if .Values.workloads.statefulSet.serviceName }}
{{- .Values.workloads.statefulSet.serviceName }}
{{- else }}
{{- printf "%s-headless" (include "helm-blueprint.fullname" .) }}
{{- end }}
{{- end }}

{{/*
Return the service name for a given service key.
Usage: {{ include "helm-blueprint.networking.serviceName" (dict "key" "http" "context" $) }}
*/}}
{{- define "helm-blueprint.networking.serviceName" -}}
{{- printf "%s-%s" (include "helm-blueprint.fullname" .context) .key }}
{{- end }}
