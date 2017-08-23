# changelog

[![travis Status](https://travis-ci.org/mh-cbon/changelog.svg?branch=master)](https://travis-ci.org/mh-cbon/changelog)[![Appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/changelog?branch=master&svg=true)](https://ci.appveyor.com/projects/mh-cbon/changelog)[![GoDoc](https://godoc.org/github.com/mh-cbon/changelog?status.svg)](http://godoc.org/github.com/mh-cbon/changelog)

Maintain a changelog easily.


This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# TOC
- [Install](#install)
  - [Glide](#glide)
  - [Chocolatey](#chocolatey)
  - [linux rpm/deb repository](#linux-rpmdeb-repository)
  - [linux rpm/deb standalone package](#linux-rpmdeb-standalone-package)
- [Usage](#usage)
  - [intermediary changelog file](#intermediary-changelog-file)
- [CLI Usage](#cli-usage)
  - [Init](#init)
  - [Prepare](#prepare)
  - [Finalize](#finalize)
- [History](#history)

# Install

Check the [release page](https://github.com/mh-cbon/changelog/releases)!

#### Glide
```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/changelog
cd $GOPATH/src/github.com/mh-cbon/changelog
git clone https://github.com/mh-cbon/changelog.git .
glide install
go install
```

#### Chocolatey
```sh
choco install changelog
```

#### linux rpm/deb repository
```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/changelog sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/changelog sh -xe
```

#### linux rpm/deb standalone package
```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/changelog sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/changelog sh -xe
```

# Usage

The workflow would be so,

- make a new repo
- commit your stuff
- run `changelog init` to generate a `change.log` file
- before a new release, run `changelog prepare`, to generate an `UNRELEASED` version
- review and edit the new `UNRELEASED` version changes in `change.log` file
- run `changelog finalize --version=x.x.x` to rename `UNRELEASED` version to its release version
- run `changelog md --out=CHANGELOG.md` to generate the new markdowned changelog file
- run `changelog debian` to get a the new version changelog to DEBIAN format

### intermediary changelog file

To work `changelog` uses an intermediary file `change.log`.

##### General overview

A `change.log` file contains a list `version` and their changes.

```
0.9.12-1

  * Initial release (Closes: #nnnn)
  * This is my first Debian package.

  - mh-cbon <mh-cbon@users.noreply.github.com>

-- Josip Rodin <joy-mg@debian.org>; Mon, 22 Mar 2010 00:37:31 +0100



0.9.12-0; distribution=unstable; urgency=low

  * Initial release (Closes: #nnnn)
  * This is my first Debian package.

  - mh-cbon <mh-cbon@users.noreply.github.com>

-- Josip Rodin <joy-mg@debian.org>; Mon, 22 Mar 2010 00:37:30 +0100
```

##### Version block

Each version is a formatted block of text such as

```
semver

  * Text of change #1
  * Text of change #2
  * ...

  - contributor #1 <mail@of.contributor.com>
  - contributor #2 <mail@of.contributor.com>

-- Packager name <mail@packager.org>; Release date
```

The most minimal version would be

```
semver

-- Packager name <mail@packager.org>; Release date
```

##### Version field

Each version starts with its version value, a valid `semver` identifier.

`semver` identifier can have additional tags, in the form of `tagname=value;`

```
semver; tag1=value; tag2=value;

-- Packager name <mail@packager.org>; Release date
```

It is allowed that a version starts by a non valid `semver` identifier,
in which case, the version is always sorted in first.

This possibility is offered to handle next `UNRELEASED` version.

```
a version name; tag1=value; tag2=value;

-- Packager name <mail@packager.org>; Release date
```

##### Version changes

Version changes immediately follow the `semver` identifier.

They start with a `space` followed by a star `*`, they can be multi-line.

The format is similar to `\s+\*\s+(.+)`

```
semver

  * This is valid
            * This is valid too, but ugly
  *This is not valid
*This is not valid either
  * This is a multiline entry
continuing here
  * This is another multiline entry
    nicer to read, leading white spaces will be trimmed
  * This is another multiline entry \
    with a backslash to get ride of the EOL

-- Packager name <mail@packager.org>; Release date
```

##### Version contributors

Version `contributors` immediately follow the list of `changes`.

They start with a `space` followed by an hyphen `-`.

The format is similar to `\s\+-\s+(.+)`.

It is not required to provide them
in the form of `name <email>`, but it is recommended.


```
semver

  - contributor #1 <mail@of.contributor.com>
  - contributor #2 <mail@of.contributor.com>

-- Packager name <mail@packager.org>; Release date
```

##### Version ender

Versions must end with a trailing line starting by a double hyphen `--`.

The trailing line provides `package author` and `release date` separated by a semicolon `;`

```
-- Packager name <mail@packager.org>; Release date
```

## CLI Usage

changelog -help
```sh
NAME:
   changelog - Changelog helper

USAGE:
   changelog <cmd> <options>

VERSION:
   0.0.0

COMMANDS:
     init       Initialize a new changelog file
     prepare    Prepare next changelog
     finalize   Take pending next changelog, apply a version on it
     test       Test to load your changelog file and report for errors or success
     export     Export the changelog using given template
     md         Export the changelog to Markdown format
     json       Export the changelog to JSON format
     debian     Export the changelog to Debian format
     rpm        Export the changelog to RPM format
     changelog  Export the changelog to CHANGELOG format
     ghrelease  Export the changelog to GHRELEASE format
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

#### Init

changelog init -help
```sh
NAME:
   changelog init - Initialize a new changelog file

USAGE:
   changelog init [command options] [arguments...]

OPTIONS:
   --author value, -a value  Package author (default: "N/A")
   --email value, -e value   Package author email
   --since value, -s value   Since which tag should the changelog be generated
```

```sh
EXAMPLE:
  changelog init
  changelog init --since=0.0.9
  changelog init --since=0.0.9 --author=mh-cbon
```

#### Prepare

changelog prepare -help
```sh
NAME:
   changelog prepare - Prepare next changelog

USAGE:
   changelog prepare [command options] [arguments...]

OPTIONS:
   --author value, -a value  Package author (default: "N/A")
   --email value, -e value   Package author email
```

```sh
EXAMPLE:
  changelog prepare
  changelog prepare --author=mh-cbon
```

#### Finalize

changelog finalize -help
```sh
NAME:
   changelog finalize - Take pending next changelog, apply a version on it

USAGE:
   changelog finalize [command options] [arguments...]

OPTIONS:
   --version value  Version revision
```

```sh
EXAMPLE:
  changelog finalize --version=0.0.2
```

##### Test

changelog test -help
```sh
NAME:
   changelog test - Test to load your changelog file and report for errors or success

USAGE:
   changelog test [arguments...]
```

```sh
EXAMPLE:
  changelog test
```

##### Export

changelog export -help
```sh
NAME:
   changelog export - Export the changelog using given template

USAGE:
   changelog export [command options] [arguments...]

OPTIONS:
   --template value, -t value  Go template
   --version value             Only given version
   --out value, -o value       Out target (default: "-")
   --guess, -g                 Automatically guess and inject name and user variable from the cwd
   --vars value                Add more variables to the template
```

```sh
EXAMPLE:
  changelog export --out=CHANGELOG.rtf --template=rtf.go
  changelog export --out=CHANGELOG.rtf --template=rtf.go --version=0.0.2
  changelog export --out=CHANGELOG.rtf --template=rtf.go --vars='{"name":"changelog"}'       # linux rocks
  changelog export --out=CHANGELOG.rtf --template=rtf.go --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### Md

changelog md -help
```sh
NAME:
   changelog md - Export the changelog to Markdown format

USAGE:
   changelog md [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --guess, -g            Automatically guess and inject name and user variable from the cwd
   --vars value           Add more variables to the template
```

```sh
EXAMPLE:
  changelog md --out=CHANGELOG.md
  changelog md --out=CHANGELOG.md --version=0.0.2
  changelog md --out=CHANGELOG.md --vars='{"name":"changelog"}'       # linux rocks
  changelog md --out=CHANGELOG.md --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### Debian

changelog debian -help
```sh
NAME:
   changelog debian - Export the changelog to Debian format

USAGE:
   changelog debian [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --guess, -g            Automatically guess and inject name and user variable from the cwd
   --vars value           Add more variables to the template
```

```sh
EXAMPLE:
  changelog debian --out=changelog
  changelog debian --out=changelog --version=0.0.2
  changelog debian --out=changelog --vars='{"name":"changelog"}'       # linux rocks
  changelog debian --out=changelog --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### RPM

changelog rpm -help
```sh
NAME:
   changelog rpm - Export the changelog to RPM format

USAGE:
   changelog rpm [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --guess, -g            Automatically guess and inject name and user variable from the cwd
   --vars value           Add more variables to the template
```

```sh
EXAMPLE:
  changelog rpm --out=changelog
  changelog rpm --out=changelog --version=0.0.2
  changelog rpm --out=changelog --vars='{"name":"changelog"}'       # linux rocks
  changelog rpm --out=changelog --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### CHANGELOG

changelog changelog -help
```sh
NAME:
   changelog changelog - Export the changelog to CHANGELOG format

USAGE:
   changelog changelog [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --guess, -g            Automatically guess and inject name and user variable from the cwd
   --vars value           Add more variables to the template
```

```sh
EXAMPLE:
  changelog changelog --out=changelog
  changelog changelog --out=changelog --version=0.0.2
  changelog changelog --out=changelog --vars='{"name":"changelog"}'       # linux rocks
  changelog changelog --out=changelog --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### GHRELEASE

changelog ghrelease -help
```sh
NAME:
   changelog ghrelease - Export the changelog to GHRELEASE format

USAGE:
   changelog ghrelease [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --guess, -g            Automatically guess and inject name and user variable from the cwd
   --vars value           Add more variables to the template
```

```sh
EXAMPLE:
  changelog ghrelease --out=changelog --version=0.0.2
```

##### Enable debug messages

To enable debug messages, just set `VERBOSE=change*`, `VERBOSE=*` before running the command.

```sh
VERBOSE=* changelog init
```

# History

[CHANGELOG](CHANGELOG.md)
