package changelog

import (
  "strings"

  "github.com/Masterminds/semver"
)

// YVersion is a marshable version
type YVersion semver.Version

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (t *YVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	tt, err := semver.NewVersion(strings.Trim(string(s), "\""))
	if err != nil {
		return err
	}
  *t = YVersion(*tt)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (t YVersion) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (t *YVersion) MarshalJSON() ([]byte, error) {
  return []byte("\""+t.String()+"\""), nil
}

func (t *YVersion) UnmarshalJSON(b []byte) error {
  n, err := semver.NewVersion(strings.Trim(string(b), "\""))
  if err != nil {
    return err
  }

  *t = YVersion(*n)
  return nil
}

// String prints a ShortTime to String.
func (t YVersion) String() string {
  v := semver.Version(t)
  return v.String()
}
