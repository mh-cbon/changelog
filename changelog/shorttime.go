package changelog

import (
  "strings"
  "time"
  "fmt"
)

// ShortTime is a marshable time format as "Mon Jan _2 2006"
type ShortTime time.Time

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (t *ShortTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	tt, err := time.Parse("Mon Jan _2 2006", strings.Trim(string(s), "\""))
	if err != nil {
		return err
	}
	*t = ShortTime(tt)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (t ShortTime) MarshalYAML() (interface{}, error) {
  stamp := fmt.Sprintf("%s", time.Time(t).Format("Mon Jan _2 2006"))
	return stamp, nil
}

func (t *ShortTime) MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("%s", time.Time(*t).Format("Mon Jan _2 2006"))
    return []byte("\""+stamp+"\""), nil
}

func (t *ShortTime) UnmarshalJSON(b []byte) error {
  tt, err := time.Parse("Mon Jan _2 2006", strings.Trim(string(b), "\""))
  if err !=nil {
    return err
  }

  l := ShortTime(tt)
  *t = l
  return nil
}

// String prints a ShortTime to String.
func (t ShortTime) String() string {
 return time.Time(t).Format("Mon Jan _2 2006")
}
