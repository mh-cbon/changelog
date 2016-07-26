package tpls

import (
	"io"
	"io/ioutil"
	"os"
	"bytes"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/mh-cbon/changelog/changelog"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

// PrintMultilines prints multiple lines with a given prefix
// replaced by space for line with index > 0
// (vertical alignment)
func PrintMultilines (lines string, prefix string) string {
  ret := ""
  for index, line := range strings.Split(lines, "\n") {
    if index==0 {
      ret += prefix + line + "\n"
    } else {
      ret += strings.Repeat(" ", len(prefix)) + line + "\n"
    }
  }
  return strings.TrimSuffix(ret, "\n")
}

// WriteTemplateTo writes changelog content
// to out target using src template file
func WriteTemplateTo(clog changelog.Changelog, partial bool, vars map[string]interface{}, src string, out string) error {
	d, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return WriteTemplateStrTo(clog, partial, vars, string(d), out)
}

// WriteTemplateStrTo write changelog content
// to out target using given template string
func WriteTemplateStrTo(clog changelog.Changelog, partial bool, vars map[string]interface{}, tplString string, out string) error {
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

  content, err := GenerateTemplateStr(clog, partial, vars, tplString)
  if err!=nil {
    return err
  }

  _, err = writer.Write([]byte(content))
	return err
}

// GenerateTemplate generates the changelog content
// using given src template file
func GenerateTemplate(clog changelog.Changelog, partial bool, vars map[string]interface{}, src string) (string, error) {
	tplString, err := ioutil.ReadFile(src)
	if err != nil {
		return "", err
	}
	return GenerateTemplateStr(clog, partial, vars, string(tplString))
}

// GenerateTemplateStr generates changelog content
// using given template string
func GenerateTemplateStr(clog changelog.Changelog, partial bool, vars map[string]interface{}, tplString string) (string, error) {

	values := make(map[string]interface{})
	values["changelog"] = clog
	values["getTagRange"] = clog.GetTagRange
	values["partial"] = partial
	values["vars"] = vars
	values["isnil"] = IsNil
	values["debianlayout"] = changelog.DateLayouts[0]
	values["rpmlayout"] = changelog.DateLayouts[4]
	values["join"] = strings.Join
	values["printMultilines"] = PrintMultilines

	t, err := template.New("it").Parse(tplString)
	if err != nil {
		return "", err
	}

  var b bytes.Buffer
	err = t.Execute(&b, values)
  if err !=nil {
    return "", err
  }
  return b.String(), nil
}

func IsNil(args *semver.Version) bool {
	return args == nil
}
