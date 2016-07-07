package changelog

import (
  "fmt"
  "errors"
  "regexp"
  "strings"
)

var packagerRegexp *regexp.Regexp
var emailRegexp *regexp.Regexp

func init () {
  packagerRegexp = regexp.MustCompile(`^\s*([^<]+)<([^>]+)>`)
  emailRegexp = regexp.MustCompile(`^\s*<([^>]+)>`)
}

type Contributors []Contributor

type Contributor struct{
  Name string
  Email string
}

func NewContributor (s string) (Contributor, error) {
  var err error
  c := Contributor{}
  if packagerRegexp.MatchString(s) {
    k1 := packagerRegexp.FindStringSubmatch(s)
    c.Name = strings.TrimSpace(k1[1])
    c.Email = strings.TrimSpace(k1[2])

  } else if emailRegexp.MatchString(s) {
    k1 := emailRegexp.FindStringSubmatch(s)
    c.Email = strings.TrimSpace(k1[1])

  } else if len(s)>0 {
    c.Name = strings.TrimSpace(s)
  }
  if len(c.Name)+len(c.Email)<1 {
    err = errors.New(fmt.Sprintf("Not a valid contributor string: %s", s))
  }
  return c, err
}

func (c Contributor) String() string {
  r := ""
  r = c.Name
  if c.Email != "" {
    r += " <" + c.Email + ">"
  }
  return strings.TrimSpace(r)
}

func (c Contributors) ContainsByEmail(email string) bool {
  for _, v := range c {
    if v.Email==email {
      return true
    }
  }
  return false
}
func (c Contributors) ContainsByName(name string) bool {
  for _, v := range c {
    if v.Name==name {
      return true
    }
  }
  return false
}
func (c Contributors) Strings() []string {
  r := make([]string, 0)
  for _, v := range c {
    r = append(r, v.String())
  }
  return r
}
