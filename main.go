package main

import (

	"github.com/mh-cbon/changelog/load"
	"github.com/urfave/cli"
)

var VERSION = "0.0.0"

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
			Action: init,
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
					Value: "",
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
			Name:   "new",
			Usage:  "Add a new version to the changelog",
			Action: newVersion,
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
					Name:  "version",
					Value: "",
					Usage: "Version name",
				},
			},
		},
	}

	app.Run(os.Args)
}

func init(c *cli.Context) error {
	name := c.String("name")
	email := c.String("email")
	author := c.String("author")
  file := "changelog.yml"

  if _, err := os.Stat(file); os.IsNotExist(err) {
    return cli.NewExitError("Changelog file exists.", 1)
  }

  c := load.Changelog{}
  c.Name = name
  c.Email = email
  c.Author = author
  err := c.Write(file)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

	return nil
}

func newVersion(c *cli.Context) error {
	

	return nil
}
