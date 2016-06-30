package tpls

var MD = `{{if eq .partial false}}# Changelog - {{.changelog.Name}}
{{end}}

{{- range $e := .changelog.Versions}}
### {{if call $.isnil $e.Version }}{{ $e.Name }}{{else}}{{ $e.Version.String }}{{end}}
{{if (gt ($e.Author|len) 0) or (gt ($e.Email|len) 0)}}
__Releaser__: {{$e.Author}}{{if gt ($e.Email|len) 0}} <{{$e.Email}}>{{end}}
{{end}}
__Date__: {{$e.Date}}
{{if gt ($e.Contributors | len) 0}}
__Contributors__: {{call $.join $e.Contributors ", "}}
{{- end}}

##### Changes
{{range $update := $e.Updates}}
- {{$update}}
{{- end}}
{{- end}}`
