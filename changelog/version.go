package changelog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
)

// Version is a struct with
// a Version field, its semver value
// a Name field, the release name
// a Date field, the release date value
// a DateLayout, the layout used to parse the input date
// Author, the release author
// Tags, the list tag associated to a release (name=value...)
// Changes, the list of changes
// Contributors, the list of contributors participating to this version
type Version struct {
	Version      *semver.Version   // semver version
	Name         string            // release name
	Date         time.Time         // date of release
	DateLayout   string            // layout of the date
	Author       Contributor       // release author
	Tags         map[string]string // deb specifics, target distributions
	Changes      []string          // list of changes
	Contributors Contributors      // list of contributors
}

// DateLayouts is a list of common date format to parse
var DateLayouts []string

func init() {
	DateLayouts = make([]string, 0)
	DateLayouts = append(DateLayouts, "Mon, 02 Jan 2006 15:04:05 -0700")
	DateLayouts = append(DateLayouts, "Mon, 02 Jan 2006 15:04:05")
	DateLayouts = append(DateLayouts, "Mon, 02 Jan 2006")
	DateLayouts = append(DateLayouts, "Mon 02 Jan 2006")
	DateLayouts = append(DateLayouts, "Mon Jan 02 2006")
	DateLayouts = append(DateLayouts, "2006-02-01")
	DateLayouts = append(DateLayouts, "2/1/2015")
	DateLayouts = append(DateLayouts, "01/02/2015")
}

// NewVersion creates a new Version struct of given version name.
// it sets the date to today by default.
func NewVersion(version string) *Version {
	v := Version{}
	v.SetTodayDate()
	vV := v.SetVersion(version)
	if vV != nil {
		v.Name = version
	}
	v.Tags = make(map[string]string)
	return &v
}

// VersionList can sort a []Version.
type VersionList []*Version

// sort utils.
func (s VersionList) Len() int {
	return len(s)
}
func (s VersionList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s VersionList) Less(i, j int) bool {
	if s[i].Version == nil {
		// true => version without version number displays on top, it is desirable for next
		return true
	}
	if s[j].Version == nil {
		// true => version without version number displays on top, it is desirable for next
		return false
	}
	v1 := semver.Version(*s[i].Version)
	v2 := semver.Version(*s[j].Version)
	return v1.GreaterThan(&v2)
}

// SetVersion interprets given string as a semver value,
// if its a valid semver, assigns it to Version property,
// otherwise it returns an error
func (v *Version) SetVersion(version string) error {
	nv, err := semver.NewVersion(version)
	if err == nil {
		v.Version = nv
	}
	return err
}

// SetDate interprets given string as a date,
// first layout to match without error is applied to Date property,
// otherwise it returns an error
func (v *Version) SetDate(date string) error {
	for _, layout := range DateLayouts {
		tt, err := time.Parse(layout, date)
		if err == nil {
			v.Date = tt
			v.DateLayout = layout
			return nil
		}
	}
	return errors.New("Failed to parse date '" + date + "'")
}

// SetTodayDate set date property to today's date.
func (v *Version) SetTodayDate() error {
	date := time.Now().Format(DateLayouts[0])
	return v.SetDate(date)
}

// GetDate returns the date string given its original layout.
func (v *Version) GetDate() string {
	return v.GetDateF(v.DateLayout)
}

// GetDateF returns the date string given a layout.
func (v *Version) GetDateF(layout string) string {
	return v.Date.Format(layout)
}

// GetName returns the version name, if empty returns its semver version.
func (v *Version) GetName() string {
	if v.Name != "" {
		return v.Name
	}
	if v.Version == nil {
		return v.Name
	}
	return v.Version.String()
}

// GetTag finds the value of tag given its name.
func (v *Version) GetTag(name string) string {
	if _, ok := v.Tags[name]; ok == false {
		return ""
	}
	return v.Tags[name]
}

// AddStrTag parses and adds given tag string with the format name=tag
func (v *Version) AddStrTag(tag string) error {
	tagRegexp := regexp.MustCompile(`\s*([^=]+)=([^=]+)`)
	tag = strings.TrimSpace(tag)
	if tagRegexp.MatchString(tag) {
		k := tagRegexp.FindStringSubmatch(tag)
		if len(k) > 1 {
			name := string(k[1])
			value := string(k[2])
			v.Tags[name] = value
			return nil
		}
	}
	return fmt.Errorf("Invalid tag '%s'", tag)
}
