package changelog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/mh-cbon/go-repo-utils/repoutils"
)

// Changelog struct contains
// a list of versions and their changes,
// a hash to the first revision of the repo.
type Changelog struct {
	Versions []*Version
	FirstRev string
}

// Load given path into the current Changelog object
func (g *Changelog) Load(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	cwd := path.Dir(filepath)
	vcs, err := repoutils.WhichVcs(cwd)
	if err == nil {
		rev, err := repoutils.GetFirstRevision(vcs, cwd)
		if err == nil {
			g.FirstRev = rev
		}
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

	versionRegexp := regexp.MustCompile(`^[^\s+].+`)
	versionEndRegexp := regexp.MustCompile(`^[-]{2}\s+(.+)`)
	changeRegexp := regexp.MustCompile(`^\s+\*\s+(.+)`)
	contributorRegexp := regexp.MustCompile(`^\s+-\s+(.+)`)
	frontspaceRegexp := regexp.MustCompile(`^(\s+)`)

	var cVersion *Version
	cChange := ""
	cVersionHasEnded := true
	lastFrontWs := 0
	var getFrontSpace = func(line []byte) int {
		k := frontspaceRegexp.FindSubmatch(line)
		return len(k[0])
	}

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
						return fmt.Errorf("%s, at line %d", err.Error(), index)
					}
				}
			}
			cVersionHasEnded = false
		} else {
			if versionEndRegexp.Match(line) {
				if len(cChange) > 0 {
					if strings.Count(cChange, "\n") <= 1 {
						cChange = strings.TrimSpace(cChange)
					}
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
						return fmt.Errorf("%s, at line %d", err.Error(), index)
					}
				} else {
					err := cVersion.SetDate(cVersion.Author.Name)
					if err == nil {
						cVersion.Author.Name = ""
					} else {
						return fmt.Errorf("%s or %s, at line %d", "Missing date", err.Error(), index)
					}
				}

				g.Versions = append(g.Versions, cVersion)
				cVersion = nil
				cVersionHasEnded = true

			} else if contributorRegexp.Match(line) && (getFrontSpace(line) <= lastFrontWs || lastFrontWs <= 0) {
				if len(cChange) > 0 {
					if strings.Count(cChange, "\n") <= 1 {
						cChange = strings.TrimSpace(cChange)
					}
					cVersion.Changes = append(cVersion.Changes, cChange)
					cChange = ""
				}
				k := contributorRegexp.FindSubmatch(line)
				c, err := NewContributor(string(k[1]))
				if err == nil {
					cVersion.Contributors = append(cVersion.Contributors, c)
				}

			} else if changeRegexp.Match(line) && (getFrontSpace(line) <= lastFrontWs || lastFrontWs <= 0) {
				if len(cChange) > 0 {
					if strings.Count(cChange, "\n") <= 1 {
						cChange = strings.TrimRight(cChange, "\n")
					}
					cVersion.Changes = append(cVersion.Changes, cChange)
				}
				// fmt.Printf("%q\n", line)
				k := changeRegexp.FindSubmatch(line)
				cChange = string(k[1])
				lastFrontWs = getFrontSpace(line)

			} else if len(cChange) > 0 {

				if cChange[len(cChange)-1:len(cChange)] == "\\" {
					cChange = strings.TrimSpace(cChange[0:len(cChange)-1]) + " " + strings.TrimSpace(string(line))

				} else if len(line) > 0 {
					x := string(line)
					if strings.TrimSpace(x) == "" {
						x = ""
					} else if lastFrontWs > 0 && strings.TrimSpace(x[:lastFrontWs]) == "" {
						x = strings.TrimRight(x[lastFrontWs:], "\n")
					}
					cChange += "\n" + x
				} else {
					cChange += "\n"
				}

			} else if len(strings.TrimSpace(string(line))) > 0 {
				return fmt.Errorf("Invalid format at line %d in %q", index, string(line))
			}
		}
	}

	if cVersionHasEnded == false {
		return fmt.Errorf("Version not closed at end of the document")
	}

	return nil
}

// FindUnreleasedVersion finds a version without name.
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

// FindVersionByName finds a version by its name.
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

// FindVersionByVersion finds a version by its version number.
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

// GetSemverVersions gets all Version with a valid semver Version field.
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

// FindMostRecentVersion finds the most recent version by comparing semver values.
func (g *Changelog) FindMostRecentVersion() *Version {
	var v *Version
	versions := g.GetSemverVersions()
	if len(versions) > 0 {
		sort.Sort(VersionList(versions))
		v = versions[0]
	}
	return v
}

// Sort ensures versions are sorted according to semver rules
func (g *Changelog) Sort() {
	sort.Sort(VersionList(g.Versions))
}

// TagRange represents a range of commit hash for a tag
type TagRange struct {
	Begin string
	End   string
}

// GetTagRange finds the commits hash range for a tag.
func (g *Changelog) GetTagRange(tag string) TagRange {
	tagRange := TagRange{}
	versions := g.GetSemverVersions()
	found := false
	for i, v := range versions {
		strV := v.GetName()
		if i == len(versions)-1 {
			tagRange.Begin = g.FirstRev
			tagRange.End = strV
		} else {
			tagRange.Begin = versions[i+1].GetName()
			tagRange.End = strV
		}
		if strV == tag {
			found = true
			break
		}
	}
	if !found {
		return TagRange{}
	}
	return tagRange
}
