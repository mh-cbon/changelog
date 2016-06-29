package tpls

import (
	"os"
	"strings"
	"text/template"

	"github.com/mh-cbon/changelog/changelog"
)

func GenerateTemplate(clog changelog.Changelog, partial bool, src string, out string) error {
	tpl, err := template.ParseFiles(src)
	if err != nil {
		return err
	}

  var writer *os.File
  if out=="-" {
    writer = os.Stdout
  } else {
  	writer, err = os.Create(out)
  	if err != nil {
  		return err
  	}
  	defer writer.Close()
  }

  values := make(map[string]interface{})
  values["changelog"] = clog
  values["partial"] = partial
  values["isnil"] = IsNil
  values["join"] = strings.Join

  template, err := tpl.New("md.go").ParseFiles(src)
  if err!=nil {
    return err
  }
	return template.Execute(writer, values)
}

func IsNil(args *changelog.YVersion) bool {
    return args==nil
}
