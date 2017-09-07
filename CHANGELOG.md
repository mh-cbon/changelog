# Changelog - changelog

### 0.0.29

__Changes__

- `changelog show`: quickly display the last version on terminal.
- `changelog rename`: quickly rename versions.
  Wihthout arguments it renames the __last version__ to `UNRELEASED`.
  Use `--version` and `--to=version|name` to control the behavior of `rename`.


















__Contributors__

- mh-cbon

Released by mh-cbon, Thu 24 Aug 2017 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.28...0.0.29#diff)
______________

### 0.0.28

__Changes__

- ci: add go 1.9, various updates
- close #5: UNRELEASED can be exported
- export: ensure the file is open with truncate flag
- close #7: the init operation does not fail on vcs issue
- changelog: 0.0.27

__Contributors__

- mh-cbon

Released by mh-cbon, Tue 22 Aug 2017 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.27...0.0.28#diff)
______________

### 0.0.27

__Changes__

- close #4: fix whitespace handling

__Contributors__

- mh-cbon

Released by mh-cbon, Wed 12 Apr 2017 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.26...0.0.27#diff)
______________

### 0.0.26

__Changes__

- add json command to export a change.log file to JSON
- add end-to-end tests
- README: make use of emd
- dependencies: update all
- Version: ensure json serialization of semver.Version type works
- bump: move script to a .version.sh file

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 24 Feb 2017 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.25...0.0.26#diff)
______________

### 0.0.25

__Changes__

- add new option -g/--guess to guess name and user variable from the cwd
- Code refactoring: applied go linters

__Contributors__

- mh-cbon

Released by mh-cbon, Wed 04 Jan 2017 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.24...0.0.25#diff)
______________

### 0.0.24

__Changes__

- travis: fix ghtoken value to build gh-pages
- appveyor: rename some variable
- appveyor: update gh token

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 30 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.23...0.0.24#diff)
______________

### 0.0.23

__Changes__

- packaging: fix travis script

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 30 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.22...0.0.23#diff)
______________

### 0.0.22

__Changes__

- packaging: add choco package, add deb/rpm repositories
- README: update usage and install section
- release: generate gh-release body from the changelog

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 30 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.21...0.0.22#diff)
______________

### 0.0.21

__Changes__

- MD layout: Fix diff url, it was using wrong repository name
- rpm: fix url of the rpm package

__Contributors__

- mh-cbon

Released by mh-cbon, Wed 27 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.20...0.0.21#diff)
______________

### 0.0.20

__Changes__

- firstrev: silently fail if the directory is not a repository,
  or if the system does not have required binaries































__Contributors__

- mh-cbon

Released by mh-cbon, Tue 26 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.19...0.0.20#diff)
______________

### 0.0.19

__Changes__

- Markdown format: add link to the github diff page
- cli: add ghrelease command to export changelog to a short md format
- multilines: improved support of multilines changes to correctly
  align them vertically with the prefix
- layouts: add new GHRELEASE template
- layouts: add test suite
- tagrange: add method to get begin...end range of a version

__Contributors__

- mh-cbon

Released by mh-cbon, Tue 26 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.18...0.0.19#diff)
______________

### 0.0.18

__Changes__

- travs: fix wrong asset name

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.17...0.0.18#diff)
______________

### 0.0.17

__Changes__

- rpm: add missing rpm.json to build the package

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.16...0.0.17#diff)
______________

### 0.0.16

__Changes__

- rpm: add missing rpm.json to build the package

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.15...0.0.16#diff)
______________

### 0.0.15

__Changes__

- rpm: add missing rpm.json to build the package

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.14...0.0.15#diff)
______________

### 0.0.14

__Changes__

- travis: fix build script
- travis: update docker file

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.13...0.0.14#diff)
______________

### 0.0.13

__Changes__

- travis: update docker file

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.12...0.0.13#diff)
______________

### 0.0.12

__Changes__

- travis: fix docker script, requires curl

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.11...0.0.12#diff)
______________

### 0.0.11

__Changes__

- build: add rpm packages

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.10...0.0.11#diff)
______________

### 0.0.10

__Changes__

- RPM: ensure release tag is at least=1
- README: update usage command lines
- export: Add changelog format to export commands
- tpls: Fix RPM template, Refactor to simplify code
- version: add method GetName to get name or version value to string
- Tags: Allow trailing semicolo in tags
- RPM: Fix rpm template

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.9...0.0.10#diff)
______________

### 0.0.9

__Changes__

- export: Add RPM export format
- main: Refactor export methods and ensure name vars is defined
- version: add new GetTag method to safely read a tag value
- debian layout: fix read urgency from the tags, not the vars
- README: update install section
- appveyor: build only on tag

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.8...0.0.9#diff)
______________

### 0.0.8

__Changes__

- appveyor: register go-msi path manually
- appveyor: update go-msi package
- appveyor: fix curl option to follow location redirect

__Contributors__

- mh-cbon

Released by mh-cbon, Tue 12 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.7...0.0.8#diff)
______________

### 0.0.7

__Changes__

- appveyor: fix curl option to follow location redirect

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 11 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.6...0.0.7#diff)
______________

### 0.0.6

__Changes__

- add msi build
- README: update install procedure

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 11 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.5...0.0.6#diff)
______________

### 0.0.5

__Changes__

- deb.json: make use of name token
- README: Update usage command line
- installer: fix url download

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 11 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.4...0.0.5#diff)
______________

### 0.0.4

__Changes__

- fix glide lock file

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 11 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.3...0.0.4#diff)
______________

### 0.0.3

__Changes__

- add magic install script
- update release key
- update install method

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 11 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.2...0.0.3#diff)
______________

### 0.0.2

__Changes__

- add deb.json file to generate debian packages

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 11 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/0.0.1...0.0.2#diff)
______________

### 0.0.1

__Changes__

- Initial release

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 11 Jul 2016 -
[see the diff](https://github.com/mh-cbon/changelog/compare/c589337667f6d64f5b2b2165290c20b8b4e7b40b...0.0.1#diff)
______________


