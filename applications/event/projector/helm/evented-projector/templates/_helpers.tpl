{{/*
Expand the name of the chart.
*/}}
{{- define "evented-projector.name" -}}
{{ $name := printf "%s-%s" .Release.Name .Chart.Name}}
{{- printf $name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "evented-projector.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "evented-projector.coordinator.port"}}
{{- .Values.coordinator.port | default 1313 -}}
{{- end }}

{{- define "evented-projector.projector.port"}}
{{- .Values.projector.port | default 1737 -}}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "evented-projector.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "evented-projector.labels" -}}
helm.sh/chart: {{ include "evented-projector.chart" . }}
{{ include "evented-projector.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "evented-projector.selectorLabels" -}}
app.kubernetes.io/name: {{ include "evented-projector.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "evented-projector.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "evented-projector.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
