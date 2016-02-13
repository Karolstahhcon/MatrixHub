{{- /*
Expand the name of the chart.
*/}}
{{- define "matrixhub.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- /*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "matrixhub.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- /*
Create chart name and version as used by the chart label.
*/}}
{{- define "matrixhub.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- /*
Common labels
*/}}
{{- define "matrixhub.labels" -}}
helm.sh/chart: {{ include "matrixhub.chart" . }}
{{ include "matrixhub.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Values.apiserver.labels }}
{{- range $key, $value := .Values.apiserver.labels }}
{{ $key }}: {{ $value | quote }}
{{- end }}
{{- end }}
{{- end }}

{{- /*
Selector labels
*/}}
{{- define "matrixhub.selectorLabels" -}}
app.kubernetes.io/name: {{ include "matrixhub.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- /*
MatrixHub API server labels
*/}}
{{- define "matrixhub.apiserver.labels" -}}
{{ include "matrixhub.labels" . }}
{{- end }}

{{- /*
MatrixHub API server pod labels
*/}}
{{- define "matrixhub.apiserver.podLabels" -}}
{{- if .Values.apiserver.podLabels }}
{{- range $key, $value := .Values.apiserver.podLabels }}
{{ $key }}: {{ $value | quote }}
{{- end }}
{{- end }}
{{ include "matrixhub.selectorLabels" . }}
{{- end }}

{{- /*
MatrixHub image
*/}}
{{- define "matrixhub.image" -}}
{{- $registry := .Values.apiserver.image.registry | default .Values.global.imageRegistry | default "" }}
{{- $repository := .Values.apiserver.image.repository }}
{{- $tag := .Values.apiserver.image.tag | default "latest" }}
{{- if $registry }}
