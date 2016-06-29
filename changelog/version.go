package changelog

import (
  "time"

  "github.com/Masterminds/semver"
)

type Version struct {
	Version        *YVersion       `json:"version,omitempty"`        // semver version
	Name           string          `json:"name,omitempty"`           // a version name
  Date           *ShortTime      `json:"date,omitempty"`           // date of release
	Author         string          `json:"author,omitempty"`         // release author
	Email          string          `json:"email,omitempty"`          // release author email
	Distribution   string          `json:"distribution,omitempty"`   // deb specifics, target distributions
	Urgency        string          `json:"urgency,omitempty"`        // deb specifics, urgency of release
  Updates        []string        `json:"updates,omitempty"`        // list of changes
  Contributors   []string        `json:"contributors,omitempty"`   // list of contributors
}

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
    return true
  }
  if s[j].Version==nil {
    return true
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
  tt, err := time.Parse("Mon Jan _2 2006", date)
  if err==nil {
    l := ShortTime(tt)
    v.Date = &l
  }
  return err
}

// Set date to DOD
func (v* Version) SetTodayDate () error {
  date := time.Now().Format("Mon Jan _2 2006")
  return v.SetDate(date)
}
