package tpls

// MD is the layout for a Markdwon export.
var MD = `{{if eq .partial false}}# Changelog - {{.vars.name}}
{{end}}

{{- range $e := .changelog.Versions}}
{{- $tagRange := call $.getTagRange $e.GetName}}
### {{$e.GetName}}

__Changes__
{{range $change := $e.Changes}}
{{call $.printMultilines $change "- "}}
{{- end}}
{{if gt ($e.Contributors | len) 0}}
__Contributors__
{{range $contributor := $e.Contributors}}
- {{$contributor.Name}}
{{- end}}
{{end}}
Released by {{$e.Author.Name}}, {{$e.Date.Format "Mon 02 Jan 2006"}} -
[see the diff](https://github.com/{{$.vars.user}}/{{$.vars.name}}/compare/{{$tagRange.Begin}}...{{$tagRange.End}}#diff)
______________
{{end}}

`

// CHANGELOG is the layout for a changelog export.
var CHANGELOG = `{{- range $index, $e := .changelog.Versions}}
{{$e.GetName}}
{{range $change := $e.Changes}}
{{call $.printMultilines $change "  * "}}
{{- end}}
{{range $contributor := $e.Contributors}}
  - {{$contributor}}
{{- end}}

-- {{$e.Author.String}}; {{$e.GetDate}}


{{end}}`

// DEBIAN is the layout for a debian export.
var DEBIAN = `{{- range $index, $e := .changelog.Versions}}
{{$.vars.name}} ({{$e.GetName}})
{{- with $urgency := $e.GetTag "urgency" }}
{{- if gt ($urgency|len) 0}}{{$urgency}};{{end}}{{- range $k,$v := $e.Tags}}{{$k}}={{$v}};{{- end}}{{end}}
{{range $change := $e.Changes}}
{{call $.printMultilines $change "  * "}}
{{- end}}

-- {{$e.Author.String}}  {{$e.GetDateF $.debianlayout}}

{{end}}`

// RPM is the layout for an rpm export.
var RPM = `{{- range $index, $e := .changelog.Versions}}
* {{$e.GetDateF $.rpmlayout}} {{$e.Author.String}} - {{$e.GetName}}{{if gt ($e.GetTag "release"|len) 0}}-{{$e.GetTag "release"}}{{else}}-1{{end}}
{{- range $change := $e.Changes}}
{{call $.printMultilines $change "- "}}
{{- end}}
{{end}}`

// GHRELEASE is the format for the
// gh release page : https://github.com/mh-cbon/go-repo-utils/releases,
// it should export only one version.
var GHRELEASE = `{{- range $e := .changelog.Versions}}
{{- range $change := $e.Changes}}
{{call $.printMultilines $change "- "}}
{{- end}}

__Contributors__ : {{call $.join $e.Contributors.Names ", "}}
{{- end}}
`
