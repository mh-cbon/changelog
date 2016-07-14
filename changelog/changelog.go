package changelog

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
)

//
type Changelog struct {
	Versions []*Version
}

// Load given path into the current Changelog object
func (g *Changelog) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return g.Parse(data)
}

// Parse and load given data into the current Changelog object.
// It will stop on first error encountered while parsing the data
func (g *Changelog) Parse(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	// @note, this is super dirty, let s improve that later.
	data = bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1)
	lines := bytes.Split(data, []byte("\n"))
	data = []byte("")

	versionRegexp := regexp.MustCompile(`^[^\s+].+`)
	versionEndRegexp := regexp.MustCompile(`^[-]{2}\s+(.+)`)
	changeRegexp := regexp.MustCompile(`^\s+\*\s+(.+)`)
	contributorRegexp := regexp.MustCompile(`^\s+-\s+(.+)`)

	var cVersion *Version
	cChange := ""
	cVersionHasEnded := true

	for index, line := range lines {
		if cVersion == nil {
			if versionRegexp.Match(line) == false {
				continue
			}
			sline := string(line)
			sline = strings.TrimSpace(sline)
			sline = strings.TrimSuffix(sline, ";")
			part := strings.Split(sline, ";")
			cVersion = NewVersion(string(part[0]))
			if len(part) > 1 {
				for _, tag := range part[1:] {
					err := cVersion.AddStrTag(string(tag))
					if err != nil {
						return errors.New(fmt.Sprintf("%s, at line %d", err.Error(), index))
					}
				}
			}
			cVersionHasEnded = false
		} else {
			if versionEndRegexp.Match(line) {
				if len(cChange) > 0 {
					cVersion.Changes = append(cVersion.Changes, cChange)
					cChange = ""
				}

				k := versionEndRegexp.FindSubmatch(line)
				part := bytes.Split(k[1], []byte(";"))

				c, err := NewContributor(string(part[0]))
				if err == nil {
					cVersion.Author = c
				}

				if len(part) > 1 {
					s := strings.TrimSpace(string(part[1]))
					err := cVersion.SetDate(s)
					if err != nil {
						return errors.New(fmt.Sprintf("%s, at line %d", err.Error(), index))
					}
				} else {
					err := cVersion.SetDate(cVersion.Author.Name)
					if err == nil {
						cVersion.Author.Name = ""
					} else {
						return errors.New(fmt.Sprintf("%s or %s, at line %d", "Missing date", err.Error(), index))
					}
				}

				g.Versions = append(g.Versions, cVersion)
				cVersion = nil
				cVersionHasEnded = true

			} else if contributorRegexp.Match(line) {
				if len(cChange) > 0 {
					cVersion.Changes = append(cVersion.Changes, cChange)
					cChange = ""
				}

				k := contributorRegexp.FindSubmatch(line)
				c, err := NewContributor(string(k[1]))
				if err == nil {
					cVersion.Contributors = append(cVersion.Contributors, c)
				}

			} else if changeRegexp.Match(line) {
				if len(cChange) > 0 {
					cVersion.Changes = append(cVersion.Changes, cChange)
					cChange = ""
				}
				k := changeRegexp.FindSubmatch(line)
				cChange = string(k[1])

			} else if len(line) > 0 && len(cChange) > 0 {
				if cChange[len(cChange)-1:len(cChange)] == "\\" {
					cChange = strings.TrimSpace(cChange[0:len(cChange)-1]) + " " + strings.TrimSpace(string(line))
				} else {
					cChange += "\n" + strings.TrimSpace(string(line))
				}

			} else if len(strings.TrimSpace(string(line))) > 0 {
				return errors.New(fmt.Sprintf("Invalid format at line %d in %q", index, string(line)))
			}
		}
	}

	if cVersionHasEnded == false {
		return errors.New(fmt.Sprintf("Version not closed at end of the document"))
	}

	return nil
}

// Find a version by its name.
func (g *Changelog) FindUnreleasedVersion() *Version {
	var v *Version
	for _, version := range g.Versions {
		if version.Version == nil {
			v = version
			break
		}
	}
	return v
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
		if s != nil && s.String() == sVersion {
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
		if s != nil {
			_, err := semver.NewVersion(version.Version.String())
			if err == nil {
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
	if len(versions) > 0 {
		sort.Sort(VersionList(versions))
		v = versions[0]
	}
	return v
}

// Ensures versions are sorted according to semver rules
func (g *Changelog) Sort() {
	sort.Sort(VersionList(g.Versions))
}
