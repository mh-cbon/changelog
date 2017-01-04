package changelog

import (
	"fmt"
	"regexp"
	"strings"
)

var packagerRegexp *regexp.Regexp
var emailRegexp *regexp.Regexp

func init() {
	packagerRegexp = regexp.MustCompile(`^\s*([^<]+)<([^>]+)>`)
	emailRegexp = regexp.MustCompile(`^\s*<([^>]+)>`)
}

// Contributors is an alias to []Contributor
type Contributors []Contributor

// Contributor is a struct representing a contributor information
type Contributor struct {
	Name  string
	Email string
}

// NewContributor parses given string into a Contributor struct,
// it returns an error if the string is not a valid contributor string format,
// name <email>
func NewContributor(s string) (Contributor, error) {
	var err error
	c := Contributor{}
	if packagerRegexp.MatchString(s) {
		k1 := packagerRegexp.FindStringSubmatch(s)
		c.Name = strings.TrimSpace(k1[1])
		c.Email = strings.TrimSpace(k1[2])

	} else if emailRegexp.MatchString(s) {
		k1 := emailRegexp.FindStringSubmatch(s)
		c.Email = strings.TrimSpace(k1[1])

	} else if len(s) > 0 {
		c.Name = strings.TrimSpace(s)
	}
	if len(c.Name)+len(c.Email) < 1 {
		err = fmt.Errorf("Not a valid contributor string: %s", s)
	}
	return c, err
}

// String returns the string representation of a contributor
func (c Contributor) String() string {
	r := c.Name
	if c.Email != "" {
		r += " <" + c.Email + ">"
	}
	return strings.TrimSpace(r)
}

// ContainsByEmail returns true if the contributor list contains given email.
func (c Contributors) ContainsByEmail(email string) bool {
	for _, v := range c {
		if v.Email == email {
			return true
		}
	}
	return false
}

// ContainsByName returns true if the contributor list contains given name.
func (c Contributors) ContainsByName(name string) bool {
	for _, v := range c {
		if v.Name == name {
			return true
		}
	}
	return false
}

// Strings returns a string slice of the contributors.
func (c Contributors) Strings() []string {
	r := make([]string, 0)
	for _, v := range c {
		r = append(r, v.String())
	}
	return r
}

// Names returns a string slice of the contributors name.
func (c Contributors) Names() []string {
	r := make([]string, 0)
	for _, v := range c {
		r = append(r, v.Name)
	}
	return r
}
