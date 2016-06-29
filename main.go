package main

import (
  "os"
  "os/exec"
  "fmt"
  "path/filepath"

	"github.com/mh-cbon/changelog/changelog"
	"github.com/mh-cbon/changelog/tpls"
	"github.com/mh-cbon/go-repo-utils/repoutils"
	"github.com/urfave/cli"
)

var VERSION = "0.0.0"

func main() {

  path, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  currentName := filepath.Base(path)

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
					Value: "",
					Usage: "Package author",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "Package author email",
				},
				cli.StringFlag{
					Name:  "name, n",
					Value: currentName,
					Usage: "Package name",
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
					Value: "",
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
			},
		},
	}

	app.Run(os.Args)
}

func initChangelog(c *cli.Context) error {
	name := c.String("name")
	email := c.String("email")
	author := c.String("author")
	since := c.String("since")
  file := "changelog.yml"

  if _, err := os.Stat(file); !os.IsNotExist(err) {
    return cli.NewExitError("Changelog file exists.", 1)
  }

  path, err := os.Getwd()
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }

  clog := changelog.Changelog{}
  clog.Name = name
  clog.Email = email
  clog.Author = author

  if c.IsSet("--since") {
    newVersion := clog.CreateVersion("next", "", "")
    err := setVersionChanges(newVersion, path, since, "")
    if err != nil {
      return cli.NewExitError(err.Error(), 1)
    }
  } else {
    tags, err := getVcsTags(path)
    if err != nil {
      return cli.NewExitError(err.Error(), 1)
    }

    for i, tag := range tags {
      since := tag
      to := ""
      if i+1<len(tags) {
        to = tags[i+1]
      }
      newVersion := clog.CreateVersion("", to, "")
      err := setVersionChanges(newVersion, path, since, to)
      if err != nil {
        return cli.NewExitError(err.Error(), 1)
      }
    }

    if len(tags)==0 {
      newVersion := clog.CreateVersion("next", "", "")
      err := setVersionChanges(newVersion, path, since, "")
      if err != nil {
        return cli.NewExitError(err.Error(), 1)
      }
    }

  }

  err = clog.Write(file)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  fmt.Println("changelog file created")

	return nil
}

func prepareNext(c *cli.Context) error {
	email := c.String("email")
	author := c.String("author")
  file := "changelog.yml"

  path, err := os.Getwd()
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }

  if _, err := os.Stat(file); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file does not exist.", 1)
  }

  clog := changelog.Changelog{}
  err = clog.Load(file)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  currentNext := clog.FindVersionByName("next")
  if currentNext!=nil {
    currentNext.SetTodayDate()
    currentNext.Updates = make([]string, 0)
  } else {
    currentNext = clog.CreateVersion("next", "", "")
  }

  mostRecent := clog.FindMostRecentVersion()
  if mostRecent==nil {
    setVersionChanges(currentNext, path, "", "")
  } else {
    setVersionChanges(currentNext, path, mostRecent.Version.String(), "")
  }

  currentNext.Author = author
  currentNext.Email = email

  err = clog.Write(file)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  fmt.Println("changelog file updated")

	return nil
}

func finalizeNext(c *cli.Context) error {
	version := c.String("version")
  file := "changelog.yml"

  if _, err := os.Stat(file); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file does not exist.", 1)
  }

  clog := changelog.Changelog{}
  err := clog.Load(file)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  currentNext := clog.FindVersionByName("next")
  if currentNext==nil {
    currentVersion := clog.FindVersionByVersion(version)
    if currentVersion==nil {
      return cli.NewExitError("No next version into this changelog", 1)
    }
    return cli.NewExitError("The version already exists and no next version was found into this changelog", 0) // desired to return 0 here.
  }

  currentNext.Name = ""
  err = currentNext.SetVersion(version)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  err = clog.Write(file)
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
  file := "changelog.yml"

  if _, err := os.Stat(file); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file does not exist.", 1)
  }

  clog := changelog.Changelog{}
  err := clog.Load(file)
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

  if version != "" {
    err = tpls.GenerateTemplate(clog, true, template, out)
  } else {
    err = tpls.GenerateTemplate(clog, false, template, out)
  }
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

	return nil
}

func exportToMd(c *cli.Context) error {
	version := c.String("version")
	out := c.String("out")
  file := "changelog.yml"

  binPath, err := getBinPath()
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }

  if _, err := os.Stat(file); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file does not exist.", 1)
  }

  clog := changelog.Changelog{}
  err = clog.Load(file)
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

  if version != "" {
    err = tpls.GenerateTemplate(clog, true, filepath.Join(binPath, "templates", "md.go"), out)
  } else {
    err = tpls.GenerateTemplate(clog, false, filepath.Join(binPath, "templates", "md.go"), out)
  }
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

	return nil
}

func setVersionChanges (version *changelog.Version, path string, since string, to string) error {

  vcs, err := repoutils.WhichVcs(path)
  if err != nil {
    return err
  }

  commits, err := repoutils.ListCommitsBetween(vcs, path, since, to)
  if err != nil {
    return err
  }

  version.Updates = make([]string, 0)
  version.Contributors = make([]string, 0)
  for _, commit := range commits {
    s := fmt.Sprintf("%s\n%s <%s> (%s)\n", commit.Message, commit.Author, commit.Email, commit.Date)
    contributor := fmt.Sprintf("%s <%s>", commit.Author, commit.Email)
    version.Updates = append(version.Updates, s)
    version.Contributors = append(version.Contributors, contributor)
  }

  return nil
}

func getVcsTags (path string) ([]string, error) {

  vcs, err := repoutils.WhichVcs(path)
  if err != nil {
    return make([]string, 0), err
  }

  return repoutils.List(vcs, path)
}

func getBinPath() (string, error) {
	var err error
	wd := ""
	if filepath.Base(os.Args[0]) == "main" { // go run ...
		wd, err = os.Getwd()
	} else {
		bin := ""
		bin, err = exec.LookPath(os.Args[0])
		if err == nil {
			wd = filepath.Dir(bin)
		}
	}
	return wd, err
}
