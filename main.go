package main

import (
  "os"
  "fmt"
  "encoding/json"

	"github.com/mh-cbon/changelog/changelog"
	"github.com/mh-cbon/changelog/tpls"
	"github.com/mh-cbon/go-repo-utils/repoutils"
	repocommit "github.com/mh-cbon/go-repo-utils/commit"
	"github.com/urfave/cli"
  "github.com/mh-cbon/verbose"
)

var VERSION = "0.0.0"
var logger = verbose.Auto()
var changelogFile = "change.log"

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
					Value: "N/A",
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
					Value: "N/A",
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
			Name:   "test",
			Usage:  "Test to load your changelog file and report for errors or success",
			Action: testFile,
			Flags: []cli.Flag{
			},
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
				cli.StringFlag{
					Name:  "vars",
					Value: "",
					Usage: "Add more variables to the template",
				},
			},
		},
		{
			Name:   "md",
			Usage:  "Export the changelog to Markdown",
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
    return cli.NewExitError("Changelog file exists.", 1)
  }

  path, err := os.Getwd()
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }

  clog := changelog.Changelog{}


  vcsTags, err := getVcsTags(path)
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }
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
      found = found || tag==since
      if found {
        proceed = true
      }
    }

    if proceed==false {
      continue
    }

    cSince := tag
    cTo := ""
    if i+1<len(tags) {
      cTo = tags[i+1]
    }
    newVersion := changelog.NewVersion(cTo)
    newVersion.Author.Email = email
    newVersion.Author.Name = author
    logger.Printf("list commits of=%q since=%q to=%q\n", cTo, cSince, cTo)
    err := setVersionChanges(newVersion, path, cSince, cTo)
    if err != nil {
      return cli.NewExitError(err.Error(), 1)
    }
    if cTo=="" {
      newVersion.Name = "UNRELEASED"
    }
    clog.Versions = append(clog.Versions, newVersion)
  }

  clog.Sort()
  vars := make(map[string]interface{})
  err = tpls.GenerateTemplateStr(clog, false, vars, tpls.CHANGELOG, changelogFile)
  if err!=nil {
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

  if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file does not exist.", 1)
  }

  clog := changelog.Changelog{}
  err = clog.Load(changelogFile)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  isNew := false
  currentNext := clog.FindVersionByName("UNRELEASED")
  if currentNext!=nil {
    currentNext.SetTodayDate()
    currentNext.Changes = make([]string, 0)
  } else {
    currentNext = changelog.NewVersion("UNRELEASED")
    isNew = true
  }

  mostRecent := clog.FindMostRecentVersion()
  if mostRecent==nil {
    setVersionChanges(currentNext, path, "", "")
  } else {
    setVersionChanges(currentNext, path, mostRecent.Version.String(), "")
  }

  currentNext.Author.Email = email
  currentNext.Author.Name = author

  if isNew {
    if len(currentNext.Changes)>0 {
      clog.Versions = append(clog.Versions, currentNext)
    } else {
      return cli.NewExitError("no changes detected", 1)
    }
  }

  clog.Sort()

  vars := make(map[string]interface{})
  err = tpls.GenerateTemplateStr(clog, false, vars, tpls.CHANGELOG, changelogFile)
  if err!=nil {
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
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  fmt.Println("changelog file is correct")

	return nil
}

func finalizeNext(c *cli.Context) error {
	version := c.String("version")

  if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file does not exist.", 1)
  }

  clog := changelog.Changelog{}
  err := clog.Load(changelogFile)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  currentNext := clog.FindVersionByName("UNRELEASED")
  if currentNext==nil {
    currentVersion := clog.FindVersionByVersion(version)
    if currentVersion==nil {
      return cli.NewExitError("No UNRELEASED version into this changelog", 1)
    }
    return cli.NewExitError("The version already exists and no UNRELEASED version was found into this changelog", 0) // desired to return 0 here.
  }

  currentNext.Name = ""
  err = currentNext.SetVersion(version)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  vars := make(map[string]interface{})
  err = tpls.GenerateTemplateStr(clog, false, vars, tpls.CHANGELOG, changelogFile)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  fmt.Println("changelog file updated")

	return nil
}

func exportChangelog(c *cli.Context) error {
	template := c.String("template")
	version := c.String("version")
	out := c.String("out")
	varsStr := c.String("vars")

  vars := make(map[string]interface{})
  if len(varsStr)>0 {
    if err := json.Unmarshal([]byte(varsStr), &vars); err != nil {
      return cli.NewExitError(fmt.Sprintf("Failed to decode vars: %s", err.Error()), 1)
    }
  }

  if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file does not exist.", 1)
  }

  clog := changelog.Changelog{}
  err := clog.Load(changelogFile)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  if version != "" {
    newVersions := make([]*changelog.Version, 0)
    v := clog.FindVersionByVersion(version)
    if v==nil {
      return cli.NewExitError("Version '"+version+"' not found.", 1)
    }
    newVersions = append(newVersions, v)
    clog.Versions = newVersions
  }

  partial := version!=""
  err = tpls.GenerateTemplate(clog, partial, vars, template, out)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

	return nil
}

func exportToMd(c *cli.Context) error {
	version := c.String("version")
	out := c.String("out")
	varsStr := c.String("vars")

  vars := make(map[string]interface{})
  if len(varsStr)>0 {
    if err := json.Unmarshal([]byte(varsStr), &vars); err != nil {
      return cli.NewExitError(fmt.Sprintf("Failed to decode vars: %s", err.Error()), 1)
    }
  }

  if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file does not exist.", 1)
  }

  clog := changelog.Changelog{}
  err := clog.Load(changelogFile)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  if version != "" {
    newVersions := make([]*changelog.Version, 0)
    v := clog.FindVersionByVersion(version)
    if v==nil {
      return cli.NewExitError("Version '"+version+"' not found.", 1)
    }
    newVersions = append(newVersions, v)
    clog.Versions = newVersions
  }

  partial := version!=""
  err = tpls.GenerateTemplateStr(clog, partial, vars, tpls.MD, out)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

	return nil
}

func setVersionChanges (version *changelog.Version, path string, since string, to string) error {

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

  for _, commit := range commits {
    s := fmt.Sprintf("%s", commit.Message)
    version.Changes = append(version.Changes, s)
    contains := version.Contributors.ContainsByEmail(commit.Email) || version.Contributors.ContainsByName(commit.Author)
    if contains==false {
      c := changelog.Contributor{}
      c.Name = commit.Author
      c.Email = commit.Email
      version.Contributors = append(version.Contributors, c)
    }
  }

  if len(commits)>0 {
    orderedCommits := repocommit.Commits(commits)
    orderedCommits.OrderByDate("DESC")
    d := orderedCommits[0].GetDate()
    if d!=nil {
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

func getVcsTags (path string) ([]string, error) {
  vcs, err := repoutils.WhichVcs(path)
  if err != nil {
    return make([]string, 0), err
  }

  tags, err := repoutils.List(vcs, path)
  if err==nil {
    tags = repoutils.SortSemverTags(tags)
  }
  return tags, err
}
