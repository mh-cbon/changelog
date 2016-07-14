package tpls

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/mh-cbon/changelog/changelog"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

func GenerateTemplate(clog changelog.Changelog, partial bool, vars map[string]interface{}, src string, out string) error {
	d, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return GenerateTemplateStr(clog, partial, vars, string(d), out)
}

func GenerateTemplateStr(clog changelog.Changelog, partial bool, vars map[string]interface{}, tplString string, out string) error {

	var err error
	var writer io.Writer
	if out == "-" {
		writer = os.Stdout
	} else {
		f, err := os.Create(out)
		if err != nil {
			return err
		}
		defer f.Close()
		writer = f
	}

	values := make(map[string]interface{})
	values["changelog"] = clog
	values["partial"] = partial
	values["vars"] = vars
	values["isnil"] = IsNil
	values["debianlayout"] = changelog.DateLayouts[0]
	values["rpmlayout"] = changelog.DateLayouts[4]
	values["join"] = strings.Join

	t, err := template.New("it").Parse(tplString)
	if err != nil {
		return err
	}
	return t.Execute(writer, values)
}

func IsNil(args *semver.Version) bool {
	return args == nil
}
