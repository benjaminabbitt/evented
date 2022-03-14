{{/*
Expand the name of the chart.
*/}}
{{- define "evented-command-handler.name" -}}
{{ $name := printf "%s-%s" .Release.Name .Chart.Name}}
{{- $name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "evented-command-handler.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "evented-command-handler.coordinator-port"}}
{{- .Values.command_handler.port | default 1313 -}}
{{- end }}

{{- define "evented-command-handler.business-logic-port"}}
{{- .Values.business_logic.port | default 1737 -}}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "evented-command-handler.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Labels
*/}}
{{- define "evented-command-handler.labels" -}}
helm.sh/chart: {{ include "evented-command-handler.chart" . }}
{{ include "evented-command-handler.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{/*
Selector labels
*/}}
{{- define "evented-command-handler.selectorLabels" -}}
app.kubernetes.io/name: {{ include "evented-command-handler.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
evented: {{ include "evented-command-handler.name" . }}
{{- end }}



{{/*
Create the name of the service account to use
*/}}
{{- define "evented-command-handler.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "evented-command-handler.fullname" .) .Values.serviceAccount.name | lower }}
{{- else }}
{{- default "default" .Values.serviceAccount.name  | lower }}
{{- end }}
{{- end }}
