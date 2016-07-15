# Changelog

Maintain a changelog easily.

# Install

__debian/ubuntu__

```sh

__deb/rpm__
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/changelog sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/changelog sh -xe
```

__windows__

Pick an msi package [here](https://github.com/mh-cbon/changelog/releases)

__go__

```sh
mkdir -p $GOPATH/src/github.com/mh-cbon
cd $GOPATH/src/github.com/mh-cbon
git clone https://github.com/mh-cbon/changelog.git
cd changelog
glide install
go install
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

#### General overview

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

#### Version block

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

#### Version field

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

#### Version changes

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

#### Version contributors

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

#### Version ender

Versions must end with a trailing line starting by a double hyphen `--`.

The trailing line provides `package author` and `release date` separated by a semicolon `;`

```
-- Packager name <mail@packager.org>; Release date
```

## CLI Usage

```sh
NAME:
   changelog - Changelog helper

USAGE:
   changelog <cmd> <options>

VERSION:
   0.0.0

COMMANDS:
     init      Initialize a new changelog file
     prepare   Prepare next changelog
     finalize  Take pending next changelog, apply a version on it
     test      Test to load your changelog file and report for errors or success
     export    Export the changelog using given template
     md        Export the changelog to Markdown format
     debian    Export the changelog to Debian format
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

#### Init

```sh
NAME:
   changelog init - Initialize a new changelog file

USAGE:
   changelog init [command options] [arguments...]

OPTIONS:
   --author value, -a value  Package author (default: "N/A")
   --email value, -e value   Package author email
   --since value, -s value   Since which tag should the changelog be generated

EXAMPLE:
  changelog init
  changelog init --since=0.0.9
  changelog init --since=0.0.9 --author=mh-cbon
```

#### Test

```sh
NAME:
   changelog test - Test to load your changelog file and report for errors or success

USAGE:
   changelog test [arguments...]

EXAMPLE:
  changelog test
```

#### Prepare

```sh
NAME:
   changelog prepare - Prepare next changelog

USAGE:
   changelog prepare [command options] [arguments...]

OPTIONS:
   --author value, -a value  Package author (default: "N/A")
   --email value, -e value   Package author email

EXAMPLE:
  changelog prepare
  changelog prepare --author=mh-cbon
```

#### Finalize

```sh
NAME:
   changelog finalize - Take pending next changelog, apply a version on it

USAGE:
   changelog finalize [command options] [arguments...]

OPTIONS:
   --version value  Version revision

EXAMPLE:
  changelog finalize --version=0.0.2
```

#### Export

```sh
NAME:
   changelog export - Export the changelog using given template

USAGE:
   changelog export [command options] [arguments...]

OPTIONS:
   --template value, -t value  Go template
   --version value             Only given version
   --out value, -o value       Out target (default: "-")
   --vars value                Add more variables to the template

EXAMPLE:
  changelog export --out=CHANGELOG.rtf --template=rtf.go
  changelog export --out=CHANGELOG.rtf --template=rtf.go --version=0.0.2
  changelog export --out=CHANGELOG.rtf --template=rtf.go --vars='{"name":"changelog"}'
```

#### Md

```sh
NAME:
   changelog md - Export the changelog to Markdown

USAGE:
   changelog md [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --vars value           Add more variables to the template

EXAMPLE:
  changelog md --out=CHANGELOG.md
  changelog md --out=CHANGELOG.md --version=0.0.2
  changelog md --out=CHANGELOG.md --vars='{"name":"changelog"}'
```

#### Debian

```sh
NAME:
   changelog debian - Export the changelog to Debian format

USAGE:
   changelog debian [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --vars value           Add more variables to the template

EXAMPLE:
  changelog debian --out=changelog
  changelog debian --out=changelog --version=0.0.2
  changelog debian --out=changelog --vars='{"name":"changelog"}'
```

#### RPM

```sh
NAME:
   changelog rpm - Export the changelog to RPM format

USAGE:
   changelog rpm [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --vars value           Add more variables to the template

EXAMPLE:
  changelog rpm --out=changelog
  changelog rpm --out=changelog --version=0.0.2
  changelog rpm --out=changelog --vars='{"name":"changelog"}'
```

#### CHANGELOG

```sh
NAME:
   changelog changelog - Export the changelog to CHANGELOG format

USAGE:
   changelog changelog [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
   --vars value           Add more variables to the template

EXAMPLE:
  changelog changelog --out=changelog
  changelog changelog --out=changelog --version=0.0.2
  changelog changelog --out=changelog --vars='{"name":"changelog"}'
```

#### Enable debug messages

To enable debug messages, just set `VERBOSE=change*`, `VERBOSE=*` before running the command.

```sh
VERBOSE=* changelog init
```

# Changelog

[CHANGELOG](CHANGELOG.md)
