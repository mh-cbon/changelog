# {{.Name}}

{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/godoc" .}}

{{pkgdoc}}

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# {{toc 5}}

# Install

{{template "gh/releases" .}}

#### Glide
{{template "glide/install" .}}

#### Chocolatey
{{template "choco/install" .}}

#### linux rpm/deb repository
{{template "linux/gh_src_repo" .}}

#### linux rpm/deb standalone package
{{template "linux/gh_pkg" .}}

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

{{exec "changelog" "-help" | color "sh"}}

#### Init

{{exec "changelog" "init" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog init
  changelog init --since=0.0.9
  changelog init --since=0.0.9 --author=mh-cbon
```

#### Prepare

{{exec "changelog" "prepare" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog prepare
  changelog prepare --author=mh-cbon
```

#### Show

{{exec "changelog" "show" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog show
```

#### Finalize

{{exec "changelog" "finalize" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog finalize --version=0.0.2
```

#### Rename

{{exec "changelog" "rename" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog rename --to UNRELEASED
  changelog rename --version 0.0.1 --to UNRELEASED
  changelog rename --version 0.0.1 --to 2.0.0
```

##### Test

{{exec "changelog" "test" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog test
```

#### Export

{{exec "changelog" "export" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog export --out=CHANGELOG.rtf --template=rtf.go
  changelog export --out=CHANGELOG.rtf --template=rtf.go --version=0.0.2
  changelog export --out=CHANGELOG.rtf --template=rtf.go --vars='{"name":"changelog"}'       # linux rocks
  changelog export --out=CHANGELOG.rtf --template=rtf.go --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### Md

{{exec "changelog" "md" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog md --out=CHANGELOG.md
  changelog md --out=CHANGELOG.md --version=0.0.2
  changelog md --out=CHANGELOG.md --vars='{"name":"changelog"}'       # linux rocks
  changelog md --out=CHANGELOG.md --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### Debian

{{exec "changelog" "debian" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog debian --out=changelog
  changelog debian --out=changelog --version=0.0.2
  changelog debian --out=changelog --vars='{"name":"changelog"}'       # linux rocks
  changelog debian --out=changelog --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### RPM

{{exec "changelog" "rpm" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog rpm --out=changelog
  changelog rpm --out=changelog --version=0.0.2
  changelog rpm --out=changelog --vars='{"name":"changelog"}'       # linux rocks
  changelog rpm --out=changelog --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### CHANGELOG

{{exec "changelog" "changelog" "-help" | color "sh"}}

```sh
EXAMPLE:
  changelog changelog --out=changelog
  changelog changelog --out=changelog --version=0.0.2
  changelog changelog --out=changelog --vars='{"name":"changelog"}'       # linux rocks
  changelog changelog --out=changelog --vars="{\"name\":\"changelog\"}"   # windows ....
```

##### GHRELEASE

{{exec "changelog" "ghrelease" "-help" | color "sh"}}

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
