package changelog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
)

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

var DateLayouts []string

func init() {
	DateLayouts = make([]string, 0)
	DateLayouts = append(DateLayouts, "Mon, 02 Jan 2006 15:04:05 -0700")
	DateLayouts = append(DateLayouts, "Mon, 02 Jan 2006 15:04:05")
	DateLayouts = append(DateLayouts, "Mon, 02 Jan 2006")
	DateLayouts = append(DateLayouts, "Mon 02 Jan 2006")
	DateLayouts = append(DateLayouts, "2006-02-01")
	DateLayouts = append(DateLayouts, "2/1/2015")
	DateLayouts = append(DateLayouts, "01/02/2015")
}

// Create a new Version.
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

// Set version value
func (v *Version) SetVersion(version string) error {
	nv, err := semver.NewVersion(version)
	if err == nil {
		v.Version = nv
	}
	return err
}

// Set date value
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

// Set date to DOD
func (v *Version) SetTodayDate() error {
	date := time.Now().Format(DateLayouts[0])
	return v.SetDate(date)
}

// Get date given its original layout.
func (v *Version) GetDate() string {
	return v.GetDateF(v.DateLayout)
}

// Get date given with a specific layout.
func (v *Version) GetDateF(layout string) string {
	return v.Date.Format(layout)
}

// Parse and add given string tag
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
	return errors.New(fmt.Sprintf("Invalid tag '%s'", tag))
}
