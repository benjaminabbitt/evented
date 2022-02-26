{{/*
Expand the name of the chart.
*/}}
{{- define "evented-query-handler.name" -}}
{{- $shortName := include "evented-query-handler.shortChartName" .}}
{{- printf "evtd-%s-%s" $shortName .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "evented-query-handler.shortChartName" -}}
{{- print .Chart.Name | trimPrefix "evented-"}}
{{- end}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "evented-query-handler.fullname" -}}
{{- $shortName := include "evented-query-handler.shortChartName" .}}
{{- printf "evtd-%s-%s" $shortName .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "evented-query-handler.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}


{{/*
Common labels
*/}}
{{- define "evented-query-handler.labels" -}}
helm.sh/chart: {{ include "evented-query-handler.chart" . }}
{{ include "evented-query-handler.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "evented-query-handler.selectorLabels" -}}
app.kubernetes.io/name: {{ include "evented-query-handler.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "evented-query-handler.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "evented-query-handler.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
