package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mh-cbon/changelog/changelog"
	"github.com/mh-cbon/changelog/tpls"
)

type TestingStub struct{}

func (t *TestingStub) Errorf(s string, a ...interface{}) {
	log.Fatalf(s+"\n", a...)
}

type TestingExiter struct{ t *testing.T }

func (t *TestingExiter) Errorf(s string, a ...interface{}) {
	panic(
		fmt.Errorf(s, a...),
	)
}

var binPath = "git_test/changelog"

func init() {

	t := &TestingStub{}
	mustFileExists(t, "main.go")

	os.RemoveAll("git_test")
	os.Mkdir("git_test", os.ModePerm)

	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	var err error
	binPath, err = filepath.Abs(binPath)
	mustNotErr(t, err)

	os.Remove(binPath)
	mustExecOk(t, makeCmd(".", "go", "build", "-o", binPath, "main.go"))

}

func initGitDir(t Errorer, dir string) {
	os.MkdirAll(dir, os.ModePerm)
	mustExecOk(t, makeCmd(dir, "git", "init"))
	mustExecOk(t, makeCmd(dir, "git", "config", "user.email", "john@doe.com"))
	mustExecOk(t, makeCmd(dir, "git", "config", "user.name", "John Doe"))
	mustExecOk(t, makeCmd(dir, "touch", "tomate"))
	mustExecOk(t, makeCmd(dir, "git", "add", "-A"))
	mustExecOk(t, makeCmd(dir, "git", "commit", "-m", "rev 1"))
}

func TestChangelog(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	mustNotErr(t, os.RemoveAll(dir))
	initGitDir(t, dir)

	//- init the changelog
	mustExecOk(tt, makeCmd(dir, binPath, "init"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))

	clog := mustGetChangelog(tt, dir)
	unreleased := mustHaveUnreleasedVersion(tt, clog)
	mustHaveNChanges(tt, unreleased, 1)
	mustHaveNContributors(tt, unreleased, 1)
	mustHaveAChange(tt, unreleased, 0, "rev 1")
	mustHaveAuthor(tt, unreleased, "John Doe", "john@doe.com")
	mustHaveAContributor(tt, unreleased, 0, "John Doe", "john@doe.com")

	mustExecOk(t, makeCmd(dir, "git", "add", "-A"))
	mustExecOk(t, makeCmd(dir, "git", "commit", "-m", "rev 2"))

	//- prepare the changelog
	mustExecOk(tt, makeCmd(dir, binPath, "prepare"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))

	clog = mustGetChangelog(tt, dir)
	unreleased = mustHaveUnreleasedVersion(tt, clog)
	mustHaveNChanges(tt, unreleased, 2)
	mustHaveAChange(tt, unreleased, 0, "rev 2")
	mustHaveAChange(tt, unreleased, 1, "rev 1")

	//- finalize the changelog
	mustExecOk(tt, makeCmd(dir, binPath, "finalize", "--version", "0.0.1"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))

	clog = mustGetChangelog(tt, dir)
	mustHaveNVersions(tt, clog, 1)
	v001 := mustHaveVersion(tt, clog, "0.0.1")
	mustHaveNChanges(tt, v001, 2)
	mustHaveAChange(tt, v001, 0, "rev 2")
	mustHaveAChange(tt, v001, 1, "rev 1")

	mustNotErr(tt, os.RemoveAll(dir))
}

func TestInitTwiceFails(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	mustNotErr(t, os.RemoveAll(dir))
	initGitDir(t, dir)

	mustExecOk(tt, makeCmd(dir, binPath, "init"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))
	mustNotExecOk(tt, makeCmd(dir, binPath, "init"))

	mustNotErr(tt, os.RemoveAll(dir))
}

func TestChangelog2(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	mustNotErr(t, os.RemoveAll(dir))
	initGitDir(t, dir)

	//- init the changelog
	mustExecOk(tt, makeCmd(dir, binPath, "init"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))

	clog := mustGetChangelog(tt, dir)
	unreleased := mustHaveUnreleasedVersion(tt, clog)
	mustHaveNChanges(tt, unreleased, 1)
	mustHaveNContributors(tt, unreleased, 1)
	mustHaveAChange(tt, unreleased, 0, "rev 1")
	mustHaveAuthor(tt, unreleased, "John Doe", "john@doe.com")
	mustHaveAContributor(tt, unreleased, 0, "John Doe", "john@doe.com")

	mustExecOk(t, makeCmd(dir, "git", "add", "-A"))
	mustExecOk(t, makeCmd(dir, "git", "commit", "-m", "rev 2"))

	//- prepare the changelog
	mustExecOk(tt, makeCmd(dir, binPath, "prepare"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))

	clog = mustGetChangelog(tt, dir)
	unreleased = mustHaveUnreleasedVersion(tt, clog)
	mustHaveNChanges(tt, unreleased, 2)
	mustHaveAChange(tt, unreleased, 0, "rev 2")
	mustHaveAChange(tt, unreleased, 1, "rev 1")

	//- finalize the changelog
	mustExecOk(tt, makeCmd(dir, binPath, "finalize", "--version", "0.1.0"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))
	mustExecOk(t, makeCmd(dir, "git", "commit", "-am", "changelog 0.1.0"))

	clog = mustGetChangelog(tt, dir)
	mustHaveNVersions(tt, clog, 1)
	v001 := mustHaveVersion(tt, clog, "0.1.0")
	mustHaveNChanges(tt, v001, 2)
	mustHaveAChange(tt, v001, 0, "rev 2")
	mustHaveAChange(tt, v001, 1, "rev 1")

	//- create branch 0.1.0
	mustExecOk(t, makeCmd(dir, "git", "tag", "0.1.0"))

	//- add a change on master
	mustExecOk(t, makeCmd(dir, "touch", "tomate-master"))
	mustExecOk(t, makeCmd(dir, "git", "add", "-A"))
	mustExecOk(t, makeCmd(dir, "git", "commit", "-m", "rev 3 master"))

	mustExecOk(tt, makeCmd(dir, binPath, "prepare"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))
	mustExecOk(tt, makeCmd(dir, binPath, "finalize", "--version", "0.1.1"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))
	mustExecOk(t, makeCmd(dir, "git", "commit", "-am", "changelog 0.1.1"))

	//- create branch 0.1.1
	mustExecOk(t, makeCmd(dir, "git", "tag", "0.1.1"))

	//- add a change on 0.1.0
	mustExecOk(t, makeCmd(dir, "git", "checkout", "0.1.0"))
	mustExecOk(t, makeCmd(dir, "touch", "tomate-0-1-0"))
	mustExecOk(t, makeCmd(dir, "git", "add", "-A"))
	mustExecOk(t, makeCmd(dir, "git", "commit", "-m", "rev 4 0.1.0"))

	mustExecOk(tt, makeCmd(dir, binPath, "prepare"))
	mustExecOk(tt, makeCmd(dir, binPath, "test"))

	clog = mustGetChangelog(tt, dir)
	mustHaveNVersions(tt, clog, 2)
	unreleased = mustHaveUnreleasedVersion(t, clog)
	mustHaveNChanges(tt, unreleased, 1)
	mustHaveAChange(tt, unreleased, 0, "rev 4 0.1.0")
	v001 = mustHaveVersion(tt, clog, "0.1.0")
	mustHaveNChanges(tt, v001, 2)
	mustHaveAChange(tt, v001, 0, "rev 2")
	mustHaveAChange(tt, v001, 1, "rev 1")

	mustExecOk(t, makeCmd(dir, "git", "commit", "-am", "changelog 0.1.0"))

	//- return on 0.1.1
	mustExecOk(t, makeCmd(dir, "git", "checkout", "0.1.1"))
	fmt.Println(mustExecOk(tt, makeCmd(dir, binPath, "json")))

	clog = mustGetChangelog(tt, dir)
	mustHaveNVersions(tt, clog, 2)
	v001 = mustHaveVersion(tt, clog, "0.1.0")
	mustHaveNChanges(tt, v001, 2)
	mustHaveAChange(tt, v001, 0, "rev 2")
	mustHaveAChange(tt, v001, 1, "rev 1")
	v011 := mustHaveVersion(tt, clog, "0.1.1")
	mustHaveNChanges(tt, v011, 1)
	mustHaveAChange(tt, v011, 0, "rev 3 master")

	mustNotErr(tt, os.RemoveAll(dir))
}

func mustGetChangelog(tt Errorer, dir string) *changelog.Changelog {
	clog := &changelog.Changelog{}
	out := mustExecOk(tt, makeCmd(dir, binPath, "json"))
	mustNotErr(tt, json.Unmarshal([]byte(out), clog))
	return clog
}

func mustHaveNVersions(t Errorer, clog *changelog.Changelog, n int) bool {
	got := len(clog.Versions)
	if got != n {
		t.Errorf(
			"Unexpected length for 'Changelog.Versions' expected=%v, got=%v",
			n,
			got,
		)
		return false
	}
	return true
}

func mustHaveVersionByName(t Errorer, clog *changelog.Changelog, name string) *changelog.Version {
	version := clog.FindVersionByName(name)
	if version == nil {
		t.Errorf(
			"Expected to find a version with name=%v",
			name,
		)
	}
	return version
}

func mustHaveVersion(t Errorer, clog *changelog.Changelog, ver string) *changelog.Version {
	version := clog.FindVersionByVersion(ver)
	if version == nil {
		t.Errorf(
			"Expected to find a version=%v",
			ver,
		)
	}
	return version
}

func mustHaveUnreleasedVersion(t Errorer, clog *changelog.Changelog) *changelog.Version {
	version := clog.FindUnreleasedVersion()
	if version == nil {
		t.Errorf(
			"Expected to find an unreleased version",
		)
	}
	return version
}

func mustHaveNChanges(t Errorer, version *changelog.Version, n int) bool {
	got := len(version.Changes)
	if got != n {
		t.Errorf(
			"Unexpected length for 'version.Changes' name=%v expected=%v, got=%v",
			version.Name,
			n,
			got,
		)
		return false
	}
	return true
}

func mustHaveNContributors(t Errorer, version *changelog.Version, n int) bool {
	got := len(version.Contributors)
	if got != n {
		t.Errorf(
			"Unexpected length for 'version.Contributors' name=%v expected=%q, got=%q",
			version.Name,
			n,
			got,
		)
		return false
	}
	return true
}

func mustHaveAChange(t Errorer, version *changelog.Version, n int, expected string) bool {
	got := version.Changes[n]
	if got != expected {
		t.Errorf(
			"Unexpected value for 'version.Changes[%v]' name=%v expected=%q, got=%q",
			n,
			version.Name,
			expected,
			got,
		)
		return false
	}
	return true
}

func mustHaveAuthor(t Errorer, version *changelog.Version, name string, email string) bool {
	got := version.Author.Name
	expected := name
	if got != expected {
		t.Errorf(
			"Unexpected value for 'version.Author.Name' name=%v expected=%q, got=%q",
			version.Name,
			expected,
			got,
		)
		return false
	}
	got = version.Author.Email
	expected = email
	if got != expected {
		t.Errorf(
			"Unexpected value for 'version.Author.Email' name=%v expected=%q, got=%q",
			version.Name,
			expected,
			got,
		)
		return false
	}
	return true
}

func mustHaveAContributor(t Errorer, version *changelog.Version, n int, name string, email string) bool {
	got := version.Contributors[n].Name
	expected := name
	if got != expected {
		t.Errorf(
			"Unexpected value for 'version.Contributors[%v].Name' name=%v expected=%q, got=%q",
			n,
			version.Name,
			expected,
			got,
		)
		return false
	}
	got = version.Contributors[n].Email
	expected = email
	if got != expected {
		t.Errorf(
			"Unexpected value for 'version.Contributors[%v].Email' name=%v expected=%q, got=%q",
			n,
			version.Name,
			expected,
			got,
		)
		return false
	}
	return true
}

type Errorer interface {
	Errorf(string, ...interface{})
}

func execOutMustContain(t Errorer, out string, s string) bool {
	if strings.Index(out, s) == -1 {
		t.Errorf("Output does not match expected to contain %q\n%v\n", s, out)
		return false
	}
	return true
}

func execOutMustNotContain(t Errorer, out string, s string) bool {
	if strings.Index(out, s) > -1 {
		t.Errorf("Output does not match expected to NOT contain %q\n%v\n", s, out)
		return false
	}
	return true
}

func mustExecOk(t Errorer, cmd *exec.Cmd) string {
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("mustExecOk failed, out=\n%v\n-------------", string(out))
	}
	mustNotErr(t, err)
	mustSucceed(t, cmd)
	return string(out)
}
func mustNotExecOk(t Errorer, cmd *exec.Cmd) string {
	out, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Println(string(out))
	}
	mustErr(t, err)
	if cmd != nil {
		mustNotSucceed(t, cmd)
	}
	return string(out)
}

func makeCmd(dir string, bin string, args ...string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	fmt.Printf("%s: %s %s\n", dir, bin, args)
	return cmd
}
func mustSucceed(t Errorer, cmd *exec.Cmd) bool {
	if cmd.ProcessState.Success() == false {
		t.Errorf("Expected success=true, got success=%t\n", false)
		return false
	}
	return true
}
func mustNotSucceed(t Errorer, cmd *exec.Cmd) bool {
	if cmd != nil && cmd.ProcessState != nil && cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%t\n", true)
		return false
	}
	return true
}
func mustErr(t Errorer, err error) bool {
	if err == nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
		return false
	}
	return true
}
func mustNotErr(t Errorer, err error) bool {
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
		return false
	}
	return true
}

func mustNotContain(t Errorer, tags []string, tag string) bool {
	if contains(tags, tag) {
		t.Errorf("Expected tags to NOT contain %q, but it WAS found in %s\n", tag, tags)
		return false
	}
	return true
}

func mustEmpty(t Errorer, tags []string) bool {
	if len(tags) > 0 {
		t.Errorf("Expected tags to be empty, but it was found %s\n", tags)
		return false
	}
	return true
}

func mustContain(t Errorer, tags []string, tag string) bool {
	if contains(tags, tag) == false {
		t.Errorf("Expected tags to contain %q, but it was not found in %s\n", tag, tags)
		return false
	}
	return true
}

func mustFileExists(t Errorer, p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Errorf("file mut exists %q", p)
		return false
	}
	return true
}
func mustWriteFile(t Errorer, p string, c string) bool {
	if err := ioutil.WriteFile(p, []byte(c), os.ModePerm); err != nil {
		t.Errorf("file not written %q", p)
		return false
	}
	return true
}

func TestFormatting(t *testing.T) {
	content := `
UNRELEASED

  * close #10: added feature to read, decode and registers the prelude data

    It is now possible to define a prelude block of ` + "`" + `yaml` + "`" + ` data in your __README__ file to
    register new data.

  * added __cat/exec/shell/color/gotest/toc__ func

    - __cat__(file string): to display the file content.
    - __exec__(bin string, args ...string): to exec a program.
    - __shell__(s string): to exec a command line on the underlying shell (it is not cross compatible).
    - __color__(color string, content string): to embed content in a block code with color.
    - __gotest__(rpkg string, run string, args ...string): exec ` + "`" + `go test <rpkg> -v -run <run> <args...>` + "`" + `.
    - __toc__(maximportance string, title string): display a TOC.

  * close #7: deprecated __file/cli__ func

    Those two functions are deprecated in flavor of their new equivalents,
    __cat/exec__.

    The new functions does not returns a triple backquuotes block code.
    They returns the response body only.
    A new function helper __color__ is a added to create a block code.

  * close #8: improved cli error output

    Before the output error was not displaying
    the command line entirely when it was too long.
    Now the error is updated to always display the command line with full length.

  * close #9: add new gotest helper func
  * close #12: add toc func
  * close#10: ensure unquoted strings are read properly
  * close #11: add shell func helper.

  - mh-cbon <mh-cbon@users.noreply.github.com>

-- mh-cbon <mh-cbon@users.noreply.github.com>; Wed, 12 Apr 2017 14:36:51 +0200

`
	s := &changelog.Changelog{}
	err := s.Parse([]byte(content))
	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	mustHaveUnreleasedVersion(t, s)

	var out bytes.Buffer
	vars := map[string]interface{}{"name": "test"}
	err = tpls.WriteTemplateStrTo(s, false, vars, tpls.MD, &out)
	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	fmt.Println(out.String())
}
