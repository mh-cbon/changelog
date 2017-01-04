package tpls

import (
	"testing"

	"github.com/mh-cbon/changelog/changelog"
)

func TestChangeLogInternalFormat(t *testing.T) {
	in := `
0.0.1

 * change #1

 - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100
`
	expected := `
0.0.1

  * change #1

  - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100


`
	s := changelog.Changelog{}
	err := s.Parse([]byte(in))
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	vars := make(map[string]interface{})
	vars["name"] = "test"

	out, err := GenerateTemplateStr(s, false, vars, CHANGELOG)
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	if expected != out {
		t.Errorf("should output=\n%q\n, got output=\n%q\n", expected, out)
	}
}

func TestMultilineFormatting(t *testing.T) {
	in := `
0.0.1

 * change #1
   with a second line
 * change #2
with another line
 * change #3
    still with another line

 - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100
`
	expected := `
0.0.1

  * change #1
    with a second line
  * change #2
    with another line
  * change #3
    still with another line

  - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100


`
	s := changelog.Changelog{}
	err := s.Parse([]byte(in))
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	vars := make(map[string]interface{})
	vars["name"] = "test"

	out, err := GenerateTemplateStr(s, false, vars, CHANGELOG)

	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	if expected != out {
		t.Errorf("should output=\n%s\n, got output=\n%s\n", expected, out)
	}
}

func TestChangeLogMDFormat(t *testing.T) {
	in := `
0.0.2

 * change #2
with another line

 - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100

0.0.1

 * change #1

 - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100
`
	expected := `# Changelog - test

### 0.0.2

__Changes__

- change #2
  with another line

__Contributors__

- contributor #1

Released by author, Mon 22 Mar 2010 -
[see the diff](https://github.com/mh-cbon/test/compare/0.0.1...0.0.2#diff)
______________

### 0.0.1

__Changes__

- change #1

__Contributors__

- contributor #1

Released by author, Mon 22 Mar 2010 -
[see the diff](https://github.com/mh-cbon/test/compare/some...0.0.1#diff)
______________


`
	s := changelog.Changelog{}
	err := s.Parse([]byte(in))
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}
	s.FirstRev = "some"

	vars := make(map[string]interface{})
	vars["name"] = "test"

	out, err := GenerateTemplateStr(s, false, vars, MD)

	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	if expected != out {
		t.Errorf("should output=\n%q\n, got output=\n%q\n", expected, out)
	}
}

func TestChangeLogDEBIANFormat(t *testing.T) {
	in := `
0.0.2

 * change #2
 with another line

 - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100

0.0.1

 * change #1

 - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100
`
	expected := `
test (0.0.2)

  * change #2
    with another line

-- author  Mon, 22 Mar 2010 00:37:30 +0100


test (0.0.1)

  * change #1

-- author  Mon, 22 Mar 2010 00:37:30 +0100

`
	s := changelog.Changelog{}
	err := s.Parse([]byte(in))
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	vars := make(map[string]interface{})
	vars["name"] = "test"

	out, err := GenerateTemplateStr(s, false, vars, DEBIAN)

	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	if expected != out {
		t.Errorf("should output=\n%q\n, got output=\n%q\n", expected, out)
	}
}

func TestChangeLogRPMFormat(t *testing.T) {
	in := `
0.0.2

 * change #2
 with another line

 - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100

0.0.1

 * change #1

 - contributor #1

-- author; Mon, 22 Mar 2010 00:37:30 +0100
`
	expected := `
* Mon Mar 22 2010 author - 0.0.2-1
- change #2
  with another line

* Mon Mar 22 2010 author - 0.0.1-1
- change #1
`
	s := changelog.Changelog{}
	err := s.Parse([]byte(in))
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	vars := make(map[string]interface{})
	vars["name"] = "test"

	out, err := GenerateTemplateStr(s, false, vars, RPM)

	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	if expected != out {
		t.Errorf("should output=\n%q\n, got output=\n%q\n", expected, out)
	}
}

func TestChangeLogGHRELEASEFormat(t *testing.T) {
	in := `
0.0.1

 * change #1
 with another line

 - contributor #1
 - contributor #2

-- author; Mon, 22 Mar 2010 00:37:30 +0100
`
	expected := `
- change #1
  with another line

__Contributors__ : contributor #1, contributor #2
`
	s := changelog.Changelog{}
	err := s.Parse([]byte(in))
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	vars := make(map[string]interface{})
	vars["name"] = "test"

	out, err := GenerateTemplateStr(s, false, vars, GHRELEASE)

	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	if expected != out {
		t.Errorf("should output=\n%q\n, got output=\n%q\n", expected, out)
	}
}
