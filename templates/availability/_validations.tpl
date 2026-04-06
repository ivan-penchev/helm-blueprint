{{/*
Validate PDB cannot have both minAvailable and maxUnavailable.
*/}}
{{- define "helm-blueprint.validation.availability.pdb" -}}
{{- if .Values.pdb.enabled }}
{{- if and .Values.pdb.minAvailable .Values.pdb.maxUnavailable }}
{{- fail "pdb: cannot set both minAvailable and maxUnavailable" }}
{{- end }}
{{- if not (or .Values.pdb.minAvailable .Values.pdb.maxUnavailable) }}
{{- fail "pdb: must set either minAvailable or maxUnavailable" }}
{{- end }}
{{- $incompatible := list "CronJob" "Job" }}
{{- if has .Values.workloadType $incompatible }}
{{- fail (printf "pdb: not applicable to workloadType %s" .Values.workloadType) }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Run all availability validations.
*/}}
{{- define "helm-blueprint.validation.availability" -}}
{{- include "helm-blueprint.validation.availability.pdb" . }}
{{- end }}
