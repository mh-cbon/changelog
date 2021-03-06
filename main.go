// Maintain a changelog easily.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/mh-cbon/changelog/changelog"
	"github.com/mh-cbon/changelog/tpls"
	repocommit "github.com/mh-cbon/go-repo-utils/commit"
	"github.com/mh-cbon/go-repo-utils/repoutils"
	"github.com/mh-cbon/verbose"
	"github.com/urfave/cli"
)

// VERSION contains the version string of the program
var VERSION = "0.0.0"
var logger = verbose.Auto()
var changelogFile = "change.log"
var notAvailable = "N/A"
var unreleased = "UNRELEASED"

func main() {

	app := cli.NewApp()
	app.Name = "changelog"
	app.Version = VERSION
	app.Usage = "Changelog helper"
	app.UsageText = "changelog <cmd> <options>"
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Initialize a new changelog file",
			Action: initChangelog,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "author, a",
					Value: notAvailable,
					Usage: "Package author",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "Package author email",
				},
				cli.StringFlag{
					Name:  "since, s",
					Value: "",
					Usage: "Since which tag should the changelog be generated",
				},
			},
		},
		{
			Name:   "prepare",
			Usage:  "Prepare next changelog",
			Action: prepareNext,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "author, a",
					Value: notAvailable,
					Usage: "Package author",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "Package author email",
				},
			},
		},
		{
			Name:   "show",
			Usage:  "Show last version",
			Action: show,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "guess, g",
					Usage: "Automatically guess and inject name and user variable from the cwd",
				},
			},
		},
		{
			Name:   "finalize",
			Usage:  "Take pending next changelog, apply a version on it",
			Action: finalizeNext,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Version revision",
				},
			},
		},
		{
			Name:   "rename",
			Usage:  "Rename a release",
			Action: rename,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Specify the version to rename",
				},
				cli.StringFlag{
					Name:  "to",
					Value: "",
					Usage: "The new name of the version",
				},
			},
		},
		{
			Name:   "test",
			Usage:  "Test to load your changelog file and report for errors or success",
			Action: testFile,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "export",
			Usage:  "Export the changelog using given template",
			Action: exportChangelog,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "template, t",
					Value: "",
					Usage: "Go template",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Only given version",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "-",
					Usage: "Out target",
				},
				cli.BoolFlag{
					Name:  "guess, g",
					Usage: "Automatically guess and inject name and user variable from the cwd",
				},
				cli.StringFlag{
					Name:  "vars",
					Value: "",
					Usage: "Add more variables to the template",
				},
			},
		},
		{
			Name:   "md",
			Usage:  "Export the changelog to Markdown format",
			Action: exportToMd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Only given version",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "-",
					Usage: "Out target",
				},
				cli.BoolFlag{
					Name:  "guess, g",
					Usage: "Automatically guess and inject name and user variable from the cwd",
				},
				cli.StringFlag{
					Name:  "vars",
					Value: "",
					Usage: "Add more variables to the template",
				},
			},
		},
		{
			Name:   "json",
			Usage:  "Export the changelog to JSON format",
			Action: exportToJSON,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Only given version",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "-",
					Usage: "Out target",
				},
			},
		},
		{
			Name:   "debian",
			Usage:  "Export the changelog to Debian format",
			Action: exportToDebian,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Only given version",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "-",
					Usage: "Out target",
				},
				cli.BoolFlag{
					Name:  "guess, g",
					Usage: "Automatically guess and inject name and user variable from the cwd",
				},
				cli.StringFlag{
					Name:  "vars",
					Value: "",
					Usage: "Add more variables to the template",
				},
			},
		},
		{
			Name:   "rpm",
			Usage:  "Export the changelog to RPM format",
			Action: exportToRpm,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Only given version",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "-",
					Usage: "Out target",
				},
				cli.BoolFlag{
					Name:  "guess, g",
					Usage: "Automatically guess and inject name and user variable from the cwd",
				},
				cli.StringFlag{
					Name:  "vars",
					Value: "",
					Usage: "Add more variables to the template",
				},
			},
		},
		{
			Name:   "changelog",
			Usage:  "Export the changelog to CHANGELOG format",
			Action: exportToChangelog,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Only given version",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "-",
					Usage: "Out target",
				},
				cli.BoolFlag{
					Name:  "guess, g",
					Usage: "Automatically guess and inject name and user variable from the cwd",
				},
				cli.StringFlag{
					Name:  "vars",
					Value: "",
					Usage: "Add more variables to the template",
				},
			},
		},
		{
			Name:   "ghrelease",
			Usage:  "Export the changelog to GHRELEASE format",
			Action: exportToGHRELEASE,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Only given version",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "-",
					Usage: "Out target",
				},
				cli.BoolFlag{
					Name:  "guess, g",
					Usage: "Automatically guess and inject name and user variable from the cwd",
				},
				cli.StringFlag{
					Name:  "vars",
					Value: "",
					Usage: "Add more variables to the template",
				},
			},
		},
	}

	app.Run(os.Args)
}

func initChangelog(c *cli.Context) error {
	email := c.String("email")
	author := c.String("author")
	since := c.String("since")

	if _, err := os.Stat(changelogFile); !os.IsNotExist(err) {
		return cli.NewExitError("Changelog file already exists.", 1)
	}

	path, err := os.Getwd()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	clog := &changelog.Changelog{}

	if vcsTags, err3 := getVcsTags(path); err3 == nil {

		tags := make([]string, 0)
		tags = append(tags, "")
		tags = append(tags, vcsTags...)

		logger.Println(since)
		logger.Println(c.IsSet("since"))

		found := false
		for i, tag := range tags {
			proceed := false

			if !c.IsSet("since") {
				proceed = true
			} else {
				found = found || tag == since
				if found {
					proceed = true
				}
			}

			if proceed == false {
				continue
			}

			cSince := tag
			cTo := ""
			if i+1 < len(tags) {
				cTo = tags[i+1]
			}
			newVersion := changelog.NewVersion(cTo)
			newVersion.Author.Email = email
			newVersion.Author.Name = author
			logger.Printf("list commits of=%q since=%q to=%q\n", cTo, cSince, cTo)
			if err2 := setVersionChanges(newVersion, path, cSince, cTo); err2 != nil && cTo != "" {
				return cli.NewExitError(err2.Error(), 1)
			}
			if newVersion.Author.Name == notAvailable && len(newVersion.Contributors) > 0 {
				newVersion.Author.Name = newVersion.Contributors[0].Name
				newVersion.Author.Email = newVersion.Contributors[0].Email
			}
			if cTo == "" {
				newVersion.Name = unreleased
			}
			clog.Versions = append(clog.Versions, newVersion)
		}
		clog.Sort()

	} else {
		// the vcs somehow is broken/notready
		// let s create an empty changelog
		newVersion := changelog.NewVersion(unreleased)
		newVersion.Author.Name = notAvailable

		clog.Versions = append(clog.Versions, newVersion)
	}

	out, err2 := os.Create(changelogFile)
	if err2 != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open file: %s", err2.Error()), 1)
	}
	vars := make(map[string]interface{})

	err = tpls.WriteTemplateStrTo(clog, false, vars, tpls.CHANGELOG, out)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error while processing the templates: %s", err.Error()), 1)
	}

	fmt.Println("changelog file created")

	return nil
}

func prepareNext(c *cli.Context) error {
	email := c.String("email")
	author := c.String("author")

	path, err := os.Getwd()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if _, err2 := os.Stat(changelogFile); os.IsNotExist(err2) {
		return cli.NewExitError("Changelog file does not exist.", 1)
	}

	clog := &changelog.Changelog{}
	err = clog.Load(changelogFile)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	isNew := false
	currentNext := clog.FindVersionByName(unreleased)
	if currentNext != nil {
		currentNext.SetTodayDate()
		currentNext.Changes = make([]string, 0)
	} else {
		currentNext = changelog.NewVersion(unreleased)
		isNew = true
	}

	mostRecent := clog.FindMostRecentVersion()
	if mostRecent == nil {
		setVersionChanges(currentNext, path, "", "")
	} else {
		setVersionChanges(currentNext, path, mostRecent.Version.String(), "")
	}

	currentNext.Author.Email = email
	currentNext.Author.Name = author

	if currentNext.Author.Name == notAvailable && len(currentNext.Contributors) > 0 {
		currentNext.Author.Name = currentNext.Contributors[0].Name
		currentNext.Author.Email = currentNext.Contributors[0].Email
	}

	if isNew {
		if len(currentNext.Changes) > 0 {
			clog.Versions = append(clog.Versions, currentNext)
		} else {
			return cli.NewExitError("no changes detected", 1)
		}
	}

	clog.Sort()

	vars := make(map[string]interface{})
	out, err2 := os.OpenFile(changelogFile, os.O_RDWR|os.O_CREATE, 0755)
	if err2 != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open file: %s", err.Error()), 1)
	}
	err = tpls.WriteTemplateStrTo(clog, false, vars, tpls.CHANGELOG, out)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("changelog file updated")

	return nil
}

func testFile(c *cli.Context) error {

	if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
		return cli.NewExitError("Changelog file does not exist.", 1)
	}

	clog := changelog.Changelog{}
	err := clog.Load(changelogFile)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("changelog file is correct")

	return nil
}

func show(c *cli.Context) error {

	guess := c.Bool("guess")
	varsStr := c.String("vars")

	if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
		return cli.NewExitError("Changelog file does not exist.", 1)
	}

	clog := &changelog.Changelog{}
	err := clog.Load(changelogFile)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	v := clog.FindUnreleasedVersion()
	if v == nil {
		v = clog.FindMostRecentVersion()
	}

	if v == nil {
		return cli.NewExitError("No version found", 1)
	}

	vars, err := computeVars(varsStr, guess)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if _, ok := vars["name"]; !ok {
		vars["name"] = ""
	}

	err2 := exportToSomeTemplate(v.GetName(), "-", vars, tpls.MD)
	if err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}
	return nil
}

func finalizeNext(c *cli.Context) error {
	version := c.String("version")

	if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
		return cli.NewExitError("Changelog file does not exist.", 1)
	}

	clog := &changelog.Changelog{}
	err := clog.Load(changelogFile)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	currentNext := clog.FindVersionByName(unreleased)
	if currentNext == nil {
		currentVersion := clog.FindVersionByVersion(version)
		if currentVersion == nil {
			return cli.NewExitError(fmt.Sprintf("No %q version into this changelog", unreleased), 1)
		}
		return cli.NewExitError(
			fmt.Sprintf("The version already exists and no %q version was found into this changelog", unreleased),
			0) // desired to return 0 here.
	}

	currentNext.Name = ""
	err = currentNext.SetVersion(version)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	vars := make(map[string]interface{})
	var out bytes.Buffer
	err = tpls.WriteTemplateStrTo(clog, false, vars, tpls.CHANGELOG, &out)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	ioutil.WriteFile(changelogFile, out.Bytes(), os.ModePerm)

	fmt.Println("changelog file updated")

	return nil
}

func rename(c *cli.Context) error {
	version := c.String("version")
	to := c.String("to")

	if to == "" {
		to = unreleased
		log.Printf("renaming to %q\n", to)
	}

	if to == version {
		return nil
	}

	if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
		return cli.NewExitError("Changelog file does not exist.", 1)
	}

	clog := &changelog.Changelog{}
	if err := clog.Load(changelogFile); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// make sure the new name does not exist
	if !(clog.FindVersionByVersion(to) == nil && clog.FindUnreleasedVersion() == nil) {
		return cli.NewExitError("version already exists: "+to, 1)
	}

	var toRename *changelog.Version
	if version == "" {
		toRename = clog.FindUnreleasedVersion()
		if toRename == nil {
			toRename = clog.FindMostRecentVersion()
		}
		log.Printf("renaming version %q\n", toRename.GetName())
	} else {
		toRename = clog.FindVersionByVersion(version)
	}

	if toRename == nil {
		return cli.NewExitError("version not found "+version, 1)
	}

	toRename.Version = nil
	toRename.Name = ""
	if _, err := semver.NewVersion(to); err == nil {
		toRename.SetVersion(to)
	} else {
		toRename.Name = to
	}

	clog.Sort()

	vars := make(map[string]interface{})
	var out bytes.Buffer
	if err := tpls.WriteTemplateStrTo(clog, false, vars, tpls.CHANGELOG, &out); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	ioutil.WriteFile(changelogFile, out.Bytes(), os.ModePerm)

	fmt.Println("changelog file updated")

	return nil
}

func exportChangelog(c *cli.Context) error {
	template := c.String("template")
	version := c.String("version")
	out := c.String("out")
	guess := c.Bool("guess")
	varsStr := c.String("vars")

	vars, err := computeVars(varsStr, guess)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	templateContent, err := ioutil.ReadFile(template)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err2 := exportToSomeTemplate(version, out, vars, string(templateContent))
	if err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	return nil
}

func exportToJSON(c *cli.Context) error {
	version := c.String("version")

	out, err := convertOut(c.String("out"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err2 := exportChangelogToJSON(version, out)
	if err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	return nil
}

func exportToMd(c *cli.Context) error {
	version := c.String("version")
	out := c.String("out")
	guess := c.Bool("guess")
	varsStr := c.String("vars")

	vars, err := computeVars(varsStr, guess)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if _, ok := vars["name"]; !ok {
		vars["name"] = ""
	}

	err2 := exportToSomeTemplate(version, out, vars, tpls.MD)
	if err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	return nil
}

func exportToDebian(c *cli.Context) error {
	version := c.String("version")
	out := c.String("out")
	guess := c.Bool("guess")
	varsStr := c.String("vars")

	vars, err := computeVars(varsStr, guess)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if _, ok := vars["name"]; !ok {
		vars["name"] = ""
	}

	err2 := exportToSomeTemplate(version, out, vars, tpls.DEBIAN)
	if err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	return nil
}

func exportToRpm(c *cli.Context) error {
	version := c.String("version")
	out := c.String("out")
	guess := c.Bool("guess")
	varsStr := c.String("vars")

	vars, err := computeVars(varsStr, guess)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err2 := exportToSomeTemplate(version, out, vars, tpls.RPM)
	if err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	return nil
}

func exportToChangelog(c *cli.Context) error {
	version := c.String("version")
	out := c.String("out")
	guess := c.Bool("guess")
	varsStr := c.String("vars")

	vars, err := computeVars(varsStr, guess)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err2 := exportToSomeTemplate(version, out, vars, tpls.CHANGELOG)
	if err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	return nil
}

func exportToGHRELEASE(c *cli.Context) error {
	version := c.String("version")
	out := c.String("out")
	guess := c.Bool("guess")
	varsStr := c.String("vars")

	vars, err := computeVars(varsStr, guess)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err2 := exportToSomeTemplate(version, out, vars, tpls.GHRELEASE)
	if err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	return nil
}

func computeVars(varsStr string, guess bool) (map[string]interface{}, error) {
	vars := make(map[string]interface{})
	if guess {
		if err := guessVars(vars); err != nil {
			return vars, cli.NewExitError(err.Error(), 1)
		}
	}
	inputVars, err2 := decodeVars(varsStr)
	if err2 != nil {
		return vars, cli.NewExitError(err2.Error(), 1)
	}
	copyVars(vars, inputVars)
	return vars, nil
}

func decodeVars(varsStr string) (map[string]interface{}, error) {
	vars := make(map[string]interface{})
	if len(varsStr) > 0 {
		if err := json.Unmarshal([]byte(varsStr), &vars); err != nil {
			return vars, fmt.Errorf("Failed to decode vars: %s (input=%v)", err.Error(), varsStr)
		}
	}
	return vars, nil
}

func copyVars(dest, src map[string]interface{}) {
	for k, v := range src {
		dest[k] = v
	}
}

func guessVars(dest map[string]interface{}) error {
	if cwd, err := os.Getwd(); err == nil {
		parts := strings.Split(cwd, string(os.PathSeparator))
		if len(parts) >= 2 {
			dest["name"] = parts[len(parts)-1]
		}
		if len(parts) >= 3 {
			dest["user"] = parts[len(parts)-2]
		}
	} else {
		return err
	}
	return nil
}

func loadChangelog(path string, version string) (*changelog.Changelog, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("Changelog file %q does not exist.", path)
	}

	clog := &changelog.Changelog{}
	err := clog.Load(changelogFile)
	if err != nil {
		return nil, err
	}

	if version != "" {
		newVersions := make([]*changelog.Version, 0)
		v := clog.FindVersionByVersion(version)
		if v == nil {
			if version == unreleased {
				v = clog.FindVersionByName(version)
				if v == nil {
					v = changelog.NewVersion(unreleased)
					v.Author.Name = notAvailable
				}
			} else {
				return nil, errors.New("Version '" + version + "' not found.")
			}
		}
		newVersions = append(newVersions, v)
		clog.Versions = newVersions
	}

	return clog, nil
}

func exportChangelogToJSON(version string, out io.Writer) error {

	clog, err := loadChangelog(changelogFile, version)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	return enc.Encode(clog)
}

func exportToSomeTemplate(version string, dst string, vars map[string]interface{}, templateContent string) error {

	clog, err := loadChangelog(changelogFile, version)
	if err != nil {
		return err
	}

	out, err := convertOut(dst)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if x, ok := out.(*os.File); ok {
		defer x.Close()
	}

	partial := version != ""
	err = tpls.WriteTemplateStrTo(clog, partial, vars, templateContent, out)
	if err != nil {
		return err
	}

	return nil
}

func setVersionChanges(version *changelog.Version, path string, since string, to string) error {

	vcs, err := repoutils.WhichVcs(path)
	if err != nil {
		logger.Printf("err=%q\n", err)
		return err
	}

	commits, err := repoutils.ListCommitsBetween(vcs, path, since, to)
	if err != nil {
		logger.Printf("err=%q\n", err)
		return err
	}

	//build changes and contributors list from the commits
	for _, commit := range commits {
		s := fmt.Sprintf("%s", commit.Message)
		version.Changes = append(version.Changes, s)
		contains := version.Contributors.ContainsByEmail(commit.Email) || version.Contributors.ContainsByName(commit.Author)
		if contains == false {
			c := changelog.Contributor{}
			c.Name = commit.Author
			c.Email = commit.Email
			version.Contributors = append(version.Contributors, c)
		}
	}

	// guess release date from the commits
	if len(commits) > 0 {
		orderedCommits := repocommit.Commits(commits)
		orderedCommits.OrderByDate("DESC")
		d := orderedCommits[0].GetDate()
		if d != nil {
			version.SetDate(d.Format(changelog.DateLayouts[0]))
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getVcsTags(path string) ([]string, error) {
	vcs, err := repoutils.WhichVcs(path)
	if err != nil {
		return make([]string, 0), err
	}

	tags, err := repoutils.List(vcs, path)
	if err == nil {
		tags = repoutils.SortSemverTags(tags)
	}
	return tags, err
}

func convertOut(out string) (io.Writer, error) {
	var writer io.Writer
	if out == "-" {
		writer = os.Stdout
	} else {
		f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return nil, err
		}
		writer = f
	}
	return writer, nil
}
