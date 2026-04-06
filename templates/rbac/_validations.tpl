{{/*
Validate roleBinding roleRef.name exists in rbac.roles.
*/}}
{{- define "helm-blueprint.validation.rbac.roleBindings" -}}
{{- range $name, $binding := .Values.rbac.roleBindings }}
{{- if not (hasKey $.Values.rbac.roles $binding.roleRef.name) }}
{{- fail (printf "rbac.roleBindings.%s: roleRef.name %q not found in rbac.roles map" $name $binding.roleRef.name) }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Validate clusterRoleBinding roleRef.name exists in rbac.clusterRoles.
*/}}
{{- define "helm-blueprint.validation.rbac.clusterRoleBindings" -}}
{{- range $name, $binding := .Values.rbac.clusterRoleBindings }}
{{- if not (hasKey $.Values.rbac.clusterRoles $binding.roleRef.name) }}
{{- fail (printf "rbac.clusterRoleBindings.%s: roleRef.name %q not found in rbac.clusterRoles map" $name $binding.roleRef.name) }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Run all RBAC validations.
*/}}
{{- define "helm-blueprint.validation.rbac" -}}
{{- include "helm-blueprint.validation.rbac.roleBindings" . }}
{{- include "helm-blueprint.validation.rbac.clusterRoleBindings" . }}
{{- end }}
