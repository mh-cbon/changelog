# Changelog

bin util to help managing a changelog

__wip__

# Install

```sh
mkdir -p $GOPATH/github.com/mh-cbon
cd $GOPATH/github.com/mh-cbon
git clone https://github.com/mh-cbon/changelog.git
cd changelog
glide install
go install
```

# Usage

The workflow would be so,

- make a new repo
- run `changelog init` to generate a changelog file
- commit your stuff
- before releasing run `changelog prepare`, it generates a `next` version
- review and edit the new `next` version in `changelog.yml`
- run `changelog finalize --version=x.x.x` to rename `next` version to its release version
- run `changelog md --out=CHANGELOG.md` to generate the new changelog file
- run `changelog deb --only=x.x.x` to get a the new version changelog and copy it somewhere else. tbd.

### intermediary changelog file

To work `changelog` uses an intermediary `changelog.yml` file.

Like this,

```yaml
name: go-repo-utils
versions:
- a_version: 0.0.7
  date: Thu Jun 30 2016
  xupdates:
  - |
    add list-commands, fix bug when detecting a non clean vcs tree
    update release script
  xcontributors:
  - mh-cbon <mh-cbon@users.noreply.github.com>
- a_version: 0.0.6
  date: Tue Jun 21 2016
  xupdates:
  - |
    Use anotated tags for git
    mh-cbon <mh-cbon@users.noreply.github.com> (Tue Jun 21 08:21:13 2016 +0200)
  xcontributors:
  - mh-cbon <mh-cbon@users.noreply.github.com>
```

few problems here

- is s using some weird names like `a_version`, `xupdates`.
Did that because i can t find a *simple* way to order `map[string]` keys when marshalling the data to yaml.
Its a bit ugly...

- It is still missing some readability, everything is very much compacted : /

## General

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
     export    Export the changelog using given template
     md        Export the changelog to Markdown
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Init

```sh
NAME:
   changelog - Initialize a new changelog file

USAGE:
   changelog [command options] [arguments...]

OPTIONS:
   --author value, -a value  Package author
   --email value, -e value   Package author email
   --name value, -n value    Package name (default: "<you pkg name>")
   --since value, -s value   Since which tag should the changelog be generated
```

## Prepare

```sh
NAME:
   changelog prepare - Prepare next changelog

USAGE:
   changelog prepare [command options] [arguments...]

OPTIONS:
   --author value, -a value  Package author
   --email value, -e value   Package author email
```

## Finalize

```sh
NAME:
   changelog finalize - Take pending next changelog, apply a version on it

USAGE:
   changelog finalize [command options] [arguments...]

OPTIONS:
   --version value  Version revision
```

## Export

```sh
NAME:
   changelog export - Export the changelog using given template

USAGE:
   changelog export [command options] [arguments...]

OPTIONS:
   --template value, -t value  Go template
   --version value             Only given version
   --out value, -o value       Out target (default: "-")
```

## Md

```sh
NAME:
   changelog md - Export the changelog to Markdown

USAGE:
   changelog md [command options] [arguments...]

OPTIONS:
   --version value        Only given version
   --out value, -o value  Out target (default: "-")
```

## Enable debug messages

To enable debug messages, just set `VERBOSE=change*`, `VERBOSE=*` before running the command.

```sh
VERBOSE=* changelog init
```

# Changelog

- 0.0.0: init
