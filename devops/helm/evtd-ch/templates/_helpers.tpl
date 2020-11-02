{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}

{{- define "commandHandler.baseName" -}}
{{- printf .Values.domain -}}
{{- end }}

{{- define "commandHandler.name" -}}
{{- printf "%s-%s" (include "commandHandler.baseName" .) .Values.defaults.commandHandlerName | trunc 63  }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "commandHandler.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" (include "commandHandler.name" .) .Release.Name | trunc 63  }}
{{- end }}
{{- end }}

{{- define "commandHandler.businessLogicName" -}}
{{- printf "%s-%s" (include "commandHandler.baseName" .) .Values.defaults.businessLogicName | trunc 63 -}}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "commandHandler.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "commandHandler.labels" -}}
helm.sh/chart: {{ include "commandHandler.chart" . }}
{{ include "commandHandler.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
evented/domain: {{.Values.domain}}
evented/name: {{.Values.name}}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "commandHandler.selectorLabels" -}}
app.kubernetes.io/name: {{ include "commandHandler.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "commandHandler.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "commandHandler.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
