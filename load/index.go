package load

import (
	"io/ioutil"
  "time"
  "fmt"
  "sort"

  "github.com/Masterminds/semver"
  "gopkg.in/yaml.v2"
)

// Config is the top-level configuration object.
type Changelog struct {
	Versions     []Version   `yaml:"versions"`
	Name         string      `yaml:"name,omitempty"`
	Author       string      `yaml:"author,omitempty"`
	Email        string      `yaml:"email,omitempty"`
}

type Version struct {
	Version        *YVersion       `yaml:"version,omitempty"`
	Name           string          `yaml:"name,omitempty"`
	Date           *ShortTime      `yaml:"date,omitempty"`
	Author         string          `yaml:"author,omitempty"`
	Email          string          `yaml:"email,omitempty"`
	Distribution   string          `yaml:"distribution,omitempty"`
	Urgency        string          `yaml:"urgency,omitempty"`
  Changes        []string        `yaml:"changes"`
}

// ShortTime is a marshable time format as "Mon Jan _2 2006"
type ShortTime time.Time

// Implements Yaml.Marshaler interface.
func (t ShortTime) MarshalYAML() ([]byte, error) {
    stamp := fmt.Sprintf("%s", time.Time(t).Format("Mon Jan _2 2006"))
    return []byte(stamp), nil
}

// Implements Yaml.Unmarshaler interface.
func (v *ShortTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
  var k string
  unmarshal(&k)

  tt, err := time.Parse("Mon Jan _2 2006", k)
  if err !=nil {
    return err
  }

  l := ShortTime(tt)
  *v = l
  return nil
}

// String prints a ShortTime to String.
func (t ShortTime) String() string {
 return time.Time(t).Format("Mon Jan _2 2006")
}



// YVersion is a marshable version
type YVersion semver.Version

// Implements Yaml.Marshaler interface.
func (t YVersion) MarshalYAML() ([]byte, error) {
  return []byte(t.String()), nil
}

// Implements Yaml.Unmarshaler interface.
func (v *YVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
  var k string
  unmarshal(&k)

  n, err := semver.NewVersion(k)
  if err != nil {
    return err
  }

  *v = YVersion(*n)
  return nil
}

// String prints a ShortTime to String.
func (t YVersion) String() string {
  v := semver.Version(t)
  return v.String()
}


// Load given path into the current Changelog object
func (g *Changelog) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return g.Parse(data)
}

// Parse and load given data into the current Changelog object
func (g *Changelog) Parse(data []byte) error {
	return yaml.Unmarshal(data, &g)
}

// Marshals to yaml then writes the Changelog instance to a file
func (g *Changelog) Write(file string) error {
  d, err := yaml.Marshal(&g)
  if err != nil {
    return err
  }
  return ioutil.WriteFile(file, d, 0644)
}

// Ensures versions are sorted according to semver rules
func (g *Changelog) Sort() {
  sort.Sort(VersionList(g.Versions))
}

// Sort implementation of []Version.
type VersionList []Version

// sort utils.
func (s VersionList) Len() int {
  return len(s)
}
func (s VersionList) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}
func (s VersionList) Less(i, j int) bool {
  if s[i].Version==nil {
    return true
  }
  if s[j].Version==nil {
    return true
  }
  v1 := semver.Version(*s[i].Version)
  v2 := semver.Version(*s[j].Version)
  return v1.GreaterThan(&v2)
}
