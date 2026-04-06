{{/*
Create the name of the service account to use
*/}}
{{- define "helm-blueprint.rbac.serviceAccountName" -}}
{{- $sa := default (dict) (default (dict) .Values.rbac).serviceAccount -}}
{{- if not (hasKey $sa "create") | ternary true $sa.create }}
{{- default (include "helm-blueprint.fullname" .) $sa.name }}
{{- else }}
{{- default "default" $sa.name }}
{{- end }}
{{- end }}
