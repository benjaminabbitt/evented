{{/*
Expand the name of the chart.
*/}}
{{- define "evented-projector.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "evented-projector.fullname" -}}
{{- printf "%s-%s" .Chart.Name .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "evented-projector.shortname"}}
{{- .Chart.Name | trimPrefix "evented-"}}
{{- end }}

{{- define "evented-projector.projector-coordinator.name" -}}
{{- $shortname := include "evented-projector.shortname" . -}}
{{- printf "evtd-pc-%s-%s" $shortname .Release.Name | lower }}
{{- end }}

{{- define "evented-projector.business-logic.name" -}}
{{- $shortname := include "evented-projector.shortname" . -}}
{{- printf "evtd-pr-%s-%s" $shortname .Release.Name | lower }}
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
