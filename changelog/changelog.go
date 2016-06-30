package changelog

import (
	"io/ioutil"
  "sort"

  "github.com/Masterminds/semver"
  "gopkg.in/yaml.v2"
)

// Config is the top-level configuration object.
type Changelog struct {
	Name         string      `yaml:"name,omitempty"`     // package name
	Author       string      `yaml:"author,omitempty"`   // default releaser author
	Email        string      `yaml:"email,omitempty"`    // default releaser email
  Versions     []*Version  `yaml:"versions"`
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
  if len(data)==0 {
    return nil
  }
	return yaml.Unmarshal(data, &g)
}

// Marshals to yaml.
func (g *Changelog) Encode() ([]byte, error) {
  return yaml.Marshal(g)
}

// Marshals to yaml then writes the Changelog instance to a file
func (g *Changelog) Write(file string) error {
  d, err := g.Encode()
  if err != nil {
    return err
  }
  return ioutil.WriteFile(file, d, 0644)
}

// Create then append a new Version on Changelog instance.
func (g *Changelog) CreateVersion(name string, version string, date string) *Version {
  v := Version{}
  v.Name = name
  if version!="" {
    v.SetVersion(version)
  }
  if date=="" {
    v.SetTodayDate()
  } else {
    v.SetDate(date)
  }
  v.Updates = make([]string, 0)
  return &v
}

// Find a version by its name.
func (g *Changelog) FindVersionByName(name string) *Version {
  var v *Version
  for _, version := range g.Versions {
    if version.Name == name {
      v = version
      break
    }
  }
  return v
}

// Find a version by its version.
func (g *Changelog) FindVersionByVersion(sVersion string) *Version {
  var v *Version
  for _, version := range g.Versions {
    s := version.Version
    if s!=nil && s.String() == sVersion {
      v = version
      break
    }
  }
  return v
}

// Get all Version with a valid semver Version field.
func (g *Changelog) GetSemverVersions() []*Version {
  ret := make([]*Version, 0)
  for _, version := range g.Versions {
    s := version.Version
    if s!=nil {
      _, err := semver.NewVersion(version.Version.String())
      if err==nil {
        ret = append(ret, version)
      }
    }
  }
  return ret
}

// Find the most recent version by comparing semver values.
func (g *Changelog) FindMostRecentVersion() *Version {
  var v *Version
  versions := g.GetSemverVersions()
  if len(versions)>0 {
    sort.Sort(VersionList(versions))
    v = versions[0]
  }
  return v
}

// Ensures versions are sorted according to semver rules
func (g *Changelog) Sort() {
  sort.Sort(VersionList(g.Versions))
}
