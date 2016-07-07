package tpls

var MD = `{{if eq .partial false}}# Changelog - {{.vars.name}}
{{end}}

{{- range $e := .changelog.Versions}}
### {{if call $.isnil $e.Version }}{{$e.Name}}{{else}}{{ $e.Version.String }}{{end}}
{{if (gt ($e.Author.Name|len) 0) or (gt ($e.Author.Email|len) 0)}}
__Releaser__: {{$e.Author.Name}}{{if gt ($e.Author.Email|len) 0}} <{{$e.Author.Email}}>{{end}}
{{end}}
__Date__: {{$e.Date.Format "Mon 02 Jan 2006"}}
{{if gt ($e.Contributors | len) 0}}
__Contributors__: {{call $.join $e.Contributors.Strings ", "}}
{{- end}}

##### Changes
{{range $change := $e.Changes}}
- {{$change}}
{{- end}}
{{- end}}
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
