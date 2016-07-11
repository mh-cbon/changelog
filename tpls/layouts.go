package tpls

var MD = `{{if eq .partial false}}# Changelog - {{.vars.name}}
{{end}}

{{- range $e := .changelog.Versions}}
### {{if call $.isnil $e.Version }}{{$e.Name}}{{else}}{{ $e.Version.String }}{{end}}

__Changes__
{{range $change := $e.Changes}}
- {{$change}}
{{end}}
{{if gt ($e.Contributors | len) 0}}
__Contributors__
{{range $contributor := $e.Contributors}}
- {{$contributor.Name}}
{{- end}}
{{end}}

Released by {{$e.Author.Name}}, {{$e.Date.Format "Mon 02 Jan 2006"}}

{{end}}

`

var CHANGELOG = `{{- range $index, $e := .changelog.Versions}}
{{- if call $.isnil $e.Version }}{{$e.Name}}{{else}}{{$e.Version.String}}{{end}}
{{range $change := $e.Changes}}
  * {{$change}}
{{- end}}
{{range $contributor := $e.Contributors}}
  - {{$contributor}}
{{- end}}

-- {{$e.Author.Name}}{{if gt ($e.Author.Email|len) 0}} <{{$e.Author.Email}}>{{end}}; {{$e.GetDate}}


{{end}}`

var DEBIAN = `{{- range $index, $e := .changelog.Versions}}
{{$.vars.name}} ({{- if call $.isnil $e.Version }}{{$e.Name}}{{else}}{{$e.Version.String}}{{end}})
{{- if gt ($.vars.urgency|len) 0}}{{$.vars.urgency}};{{end}}{{- range $k,$v := $e.Tags}}{{$k}}={{$v}};{{- end}}
{{range $change := $e.Changes}}
  * {{$change}}
{{- end}}

-- {{$e.Author.Name}}{{if gt ($e.Author.Email|len) 0}} <{{$e.Author.Email}}>{{end}}  {{$e.GetDateF $.debianlayout}}

{{end}}`
