package changelog

import (
  "time"

  "github.com/Masterminds/semver"
)

type Version struct {
	Version        *YVersion       `yaml:"a_version,omitempty"`      // semver version
	Name           string          `yaml:"aname,omitempty"`          // a version name
  Date           *ShortTime      `yaml:"date,omitempty"`           // date of release
	Author         string          `yaml:"author,omitempty"`         // release author
	Email          string          `yaml:"email,omitempty"`          // release author email
	Distribution   string          `yaml:"distribution,omitempty"`   // deb specifics, target distributions
	Urgency        string          `yaml:"urgency,omitempty"`        // deb specifics, urgency of release
  Updates        []string        `yaml:"xupdates,omitempty"`       // list of changes
  Contributors   []string        `yaml:"xcontributors,omitempty"`  // list of contributors
}

var DateLayout = "Mon Jan _2 2006"

// Sort implementation of []Version.
type VersionList []*Version

// sort utils.
func (s VersionList) Len() int {
  return len(s)
}
func (s VersionList) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}
func (s VersionList) Less(i, j int) bool {
  if s[i].Version==nil {
    // true => version without version number displays on top, it is desirable for next
    return true
  }
  if s[j].Version==nil {
     // true => version without version number displays on top, it is desirable for next
    return false
  }
  v1 := semver.Version(*s[i].Version)
  v2 := semver.Version(*s[j].Version)
  return v1.GreaterThan(&v2)
}

// Set version value
func (v* Version) SetVersion (version string) error {
  nv, err := semver.NewVersion(version)
  if err==nil {
    yv := YVersion(*nv)
    v.Version = &yv
  }
  return err
}

// Set date value
func (v* Version) SetDate (date string) error {
  tt, err := time.Parse(DateLayout, date)
  if err==nil {
    l := ShortTime(tt)
    v.Date = &l
  }
  return err
}

// Set date to DOD
func (v* Version) SetTodayDate () error {
  date := time.Now().Format(DateLayout)
  return v.SetDate(date)
}
