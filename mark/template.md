# NAME

{{ summary . }}

# USAGE

{{ usage . | indent 4 }}

{{if .Cmds -}}
# COMMANDS

{{ commands . | indent 4 }}

{{ end -}}
{{- if .Long -}}
# DESCRIPTION

{{ long . }}

{{ end }}

