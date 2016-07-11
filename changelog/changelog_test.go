package changelog

import (
	"testing"
)

func TestParseEmptyFile(t *testing.T) {
	s := Changelog{}
	err := s.Parse([]byte(""))

	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}
	if len(s.Versions) > 0 {
		t.Errorf("should len(s.Versions)==0, got len(s.Versions)==%d\n", len(s.Versions))
	}
}

func TestParseFile1(t *testing.T) {
	content := `
0.0.1
-- author; Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}
	if len(s.Versions) == 0 {
		t.Errorf("should len(s.Versions)>'0', got len(s.Versions)=%q\n", len(s.Versions))
	}
	v := s.Versions[0]
	if v.Name != "" {
		t.Errorf("should s.Name='', got s.Name=%q\n", v.Name)
	}
	if v.Version == nil {
		t.Errorf("should s.Version!='nil', got s.Version=%q\n", v.Version)
	} else if v.Version.String() != "0.0.1" {
		t.Errorf("should s.Version.String()=='0.0.1', got s.Version.String()=%q\n", v.Version.String())
	}
	if v.Author.Name != "author" {
		t.Errorf("should s.Author.Name='author', got s.Author.Name=%q\n", v.Author.Name)
	}
	if v.Author.Email != "" {
		t.Errorf("should s.Author.Email='', got s.Author.Email=%q\n", v.Author.Email)
	}
}

func TestUnclosedVersion(t *testing.T) {
	content := `
0.0.1
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err == nil {
		t.Errorf("should err!=nil, got err=%q\n", err)
	}
	if len(s.Versions) > 0 {
		t.Errorf("should len(s.Versions)='0', got len(s.Versions)=%q\n", len(s.Versions))
	}
}

func TestVersionEndAuthorEmailDate(t *testing.T) {
	content := `
0.0.1
-- author <email>; Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if v.Author.Name != "author" {
		t.Errorf("should s.Author.Name='author', got s.Author.Name=%q\n", v.Author.Name)
	}
	if v.Author.Email != "email" {
		t.Errorf("should s.Author.Email='email', got s.Author.Email=%q\n", v.Author.Email)
	}
	e := "Mon, 22 Mar 2010 00:37:30 +0100"
	if v.Date.Format(DateLayouts[0]) != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.Date.Format(e))
	}
}

func TestVersionEndAuthorDate(t *testing.T) {
	content := `
0.0.1
-- author; Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if v.Author.Name != "author" {
		t.Errorf("should s.Author.Name='author', got s.Author.Name=%q\n", v.Author.Name)
	}
	if v.Author.Email != "" {
		t.Errorf("should s.Author.Email='', got s.Author.Email=%q\n", v.Author.Email)
	}
	e := "Mon, 22 Mar 2010 00:37:30 +0100"
	if v.Date.Format(DateLayouts[0]) != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.Date.Format(e))
	}
}

func TestVersionEndEmailDate(t *testing.T) {
	content := `
0.0.1
-- <email>; Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if v.Author.Name != "" {
		t.Errorf("should s.Author.Name='', got s.Author.Name=%q\n", v.Author.Name)
	}
	if v.Author.Email != "email" {
		t.Errorf("should s.Author.Email='email', got s.Author.Email=%q\n", v.Author.Email)
	}
	e := "Mon, 22 Mar 2010 00:37:30 +0100"
	if v.Date.Format(DateLayouts[0]) != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.Date.Format(e))
	}
}

func TestVersionEndMissingDate(t *testing.T) {
	content := `
0.0.1
-- author
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err == nil {
		t.Errorf("should err!=nil, got err=%q\n", err)
	}
}

func TestVersionEndOnlyDate(t *testing.T) {
	content := `
0.0.1
-- Mon, 22 Mar 2010 00:37:30 +0100
`

	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if v.Author.Name != "" {
		t.Errorf("should s.Author.Name='', got s.Author.Name=%q\n", v.Author.Name)
	}
	if v.Author.Email != "" {
		t.Errorf("should s.Author.Email='', got s.Author.Email=%q\n", v.Author.Email)
	}
	e := "Mon, 22 Mar 2010 00:37:30 +0100"
	if v.Date.Format(DateLayouts[0]) != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.Date.Format(e))
	}
}

func TestWrongVersionAsAName(t *testing.T) {
	content := `
UNRELEASED
  - Contributor 1
  - Contributor 2
  * Change 1
  * Change 2
-- Mon, 22 Mar 2010 00:37:30 +0100
`

	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if v.Name != "UNRELEASED" {
		t.Errorf("should v.Name='UNRELEASED', got v.Name=%q\n", v.Name)
	}
}

func TestVersionTags(t *testing.T) {
	content := `
UNRELEASED; tag1=v1; tag2=v2
-- Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if val, ok := v.Tags["tag1"]; ok == false {
		t.Errorf("should 'tag1' in v.Tags, got %q\n", ok)
	} else if val != "v1" {
		t.Errorf("should v.Tags['tag1']='v1', got v.Tags['tag1']=%q\n", val)
	}
	if val, ok := v.Tags["tag2"]; ok == false {
		t.Errorf("should 'tag2' in v.Tags, got %q\n", ok)
	} else if val != "v2" {
		t.Errorf("should v.Tags['tag2']='v2', got v.Tags['tag2']=%q\n", val)
	}
}

func TestMultipleVersions(t *testing.T) {
	content := `
0.0.1
-- Mon, 22 Mar 2010 00:37:30 +0100
0.0.2
-- Mon, 22 Mar 2010 00:37:30 +0100
`

	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	if len(s.Versions) != 2 {
		t.Errorf("should len(s.Versions)='2', got len(s.Versions)=%q\n", len(s.Versions))
	}
	v := s.Versions[0]
	if v.Name != "" {
		t.Errorf("should s.Name='', got s.Name=%q\n", v.Name)
	}
	if v.Version == nil {
		t.Errorf("should s.Version!='nil', got s.Version=%q\n", v.Version)
	} else if v.Version.String() != "0.0.1" {
		t.Errorf("should s.Version.String()=='0.0.1', got s.Version.String()=%q\n", v.Version.String())
	}
	v = s.Versions[1]
	if v.Name != "" {
		t.Errorf("should s.Name='', got s.Name=%q\n", v.Name)
	}
	if v.Version == nil {
		t.Errorf("should s.Version!='nil', got s.Version=%q\n", v.Version)
	} else if v.Version.String() != "0.0.2" {
		t.Errorf("should s.Version.String()=='0.0.2', got s.Version.String()=%q\n", v.Version.String())
	}
}

func TestChanges(t *testing.T) {
	content := `
0.0.1
 * Change 1
 * Change 2
-- Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if len(v.Changes) != 2 {
		t.Errorf("should len(v.Changes)='2', got len(v.Changes)=%q\n", len(v.Changes))
	}
	if v.Changes[0] != "Change 1" {
		t.Errorf("should v.Changes[0]='Change 1', got v.Changes[0]=%q\n", v.Changes[0])
	}
	if v.Changes[1] != "Change 2" {
		t.Errorf("should v.Changes[1]='Change 2', got v.Changes[1]=%q\n", v.Changes[1])
	}
}

func TestChangesMultiline(t *testing.T) {
	content := `
0.0.1
 * Change 0
line 2
 * Change 1
   line 2
 * Change 2 \
   line 2
 * Change 3 \
line 2
-- Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if len(v.Changes) != 4 {
		t.Errorf("should len(v.Changes)='4', got len(v.Changes)=%q\n", len(v.Changes))
	}
	if v.Changes[0] != "Change 0\nline 2" {
		t.Errorf("should v.Changes[0]='Change 0\nline 2', got v.Changes[0]=%q\n", v.Changes[0])
	}
	if v.Changes[1] != "Change 1\nline 2" {
		t.Errorf("should v.Changes[1]='Change 1\nline 2', got v.Changes[1]=%q\n", v.Changes[1])
	}
	if v.Changes[2] != "Change 2 line 2" {
		t.Errorf("should v.Changes[2]='Change 2 line 2', got v.Changes[2]=%q\n", v.Changes[2])
	}
	if v.Changes[3] != "Change 3 line 2" {
		t.Errorf("should v.Changes[3]='Change 3 line 2', got v.Changes[3]=%q\n", v.Changes[3])
	}
}

func TestContributors(t *testing.T) {
	content := `
0.0.1
 - Contributor 1
 - Contributor 2
-- Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if len(v.Contributors) != 2 {
		t.Errorf("should len(v.Contributors)='2', got len(v.Contributors)=%q\n", len(v.Changes))
	}
	if v.Contributors[0].Name != "Contributor 1" {
		t.Errorf("should v.Contributors[0].Name='Contributor 1', got v.Contributors[0].Name=%q\n", v.Contributors[0].Name)
	}
	if v.Contributors[0].Email != "" {
		t.Errorf("should v.Contributors[0].Email='', got v.Contributors[0].Email=%q\n", v.Contributors[0].Email)
	}
	if v.Contributors[1].Name != "Contributor 2" {
		t.Errorf("should v.Contributors[1].Name='Contributor 2', got v.Contributors[1].Name=%q\n", v.Contributors[1].Name)
	}
	if v.Contributors[1].Email != "" {
		t.Errorf("should v.Contributors[1].Email='', got v.Contributors[1].Email=%q\n", v.Contributors[1].Email)
	}
}

func TestRandomOk1(t *testing.T) {
	content := `
0.0.1
 * Change 1
  * Change 2
 - Contributor 1
 - Contributor 2
-- Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if len(v.Contributors) != 2 {
		t.Errorf("should len(v.Contributors)='2', got len(v.Contributors)=%q\n", len(v.Changes))
	}
	if v.Contributors[0].Name != "Contributor 1" {
		t.Errorf("should v.Contributors[0].Name='Contributor 1', got v.Contributors[0].Name=%q\n", v.Contributors[0].Name)
	}
	if v.Contributors[0].Email != "" {
		t.Errorf("should v.Contributors[0].Email='', got v.Contributors[0].Email=%q\n", v.Contributors[0].Email)
	}
	if v.Contributors[1].Name != "Contributor 2" {
		t.Errorf("should v.Contributors[1].Name='Contributor 2', got v.Contributors[1].Name=%q\n", v.Contributors[1].Name)
	}
	if v.Contributors[1].Email != "" {
		t.Errorf("should v.Contributors[1].Email='', got v.Contributors[1].Email=%q\n", v.Contributors[1].Email)
	}
	if len(v.Changes) != 2 {
		t.Errorf("should len(v.Changes)='2', got len(v.Changes)=%q\n", len(v.Changes))
	}
	if v.Changes[0] != "Change 1" {
		t.Errorf("should v.Changes[0]='Change 1', got v.Changes[0]=%q\n", v.Changes[0])
	}
	if v.Changes[1] != "Change 2" {
		t.Errorf("should v.Changes[1]='Change 2', got v.Changes[1]=%q\n", v.Changes[1])
	}
}

func TestRandomOk2(t *testing.T) {
	content := `
0.0.1
  - Contributor 1
  - Contributor 2
 * Change 1
  * Change 2
-- Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}
	v := s.Versions[0]
	if len(v.Contributors) != 2 {
		t.Errorf("should len(v.Contributors)='2', got len(v.Contributors)=%q\n", len(v.Changes))
	}
	if v.Contributors[0].Name != "Contributor 1" {
		t.Errorf("should v.Contributors[0].Name='Contributor 1', got v.Contributors[0].Name=%q\n", v.Contributors[0].Name)
	}
	if v.Contributors[0].Email != "" {
		t.Errorf("should v.Contributors[0].Email='', got v.Contributors[0].Email=%q\n", v.Contributors[0].Email)
	}
	if v.Contributors[1].Name != "Contributor 2" {
		t.Errorf("should v.Contributors[1].Name='Contributor 2', got v.Contributors[1].Name=%q\n", v.Contributors[1].Name)
	}
	if v.Contributors[1].Email != "" {
		t.Errorf("should v.Contributors[1].Email='', got v.Contributors[1].Email=%q\n", v.Contributors[1].Email)
	}
	if len(v.Changes) != 2 {
		t.Errorf("should len(v.Changes)='2', got len(v.Changes)=%q\n", len(v.Changes))
	}
	if v.Changes[0] != "Change 1" {
		t.Errorf("should v.Changes[0]='Change 1', got v.Changes[0]=%q\n", v.Changes[0])
	}
	if v.Changes[1] != "Change 2" {
		t.Errorf("should v.Changes[1]='Change 2', got v.Changes[1]=%q\n", v.Changes[1])
	}
}

func TestPickyContributorMustHaveFrontSpace(t *testing.T) {
	content := `
0.0.1
  - Contributor 1
- Contributor 2
 * Change 1
  * Change 2
-- Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err == nil {
		t.Errorf("should err!=nil, got err=%q\n", err)
	}
}

func TestPickyChangeMustHaveFrontSpace(t *testing.T) {
	content := `
0.0.1
  - Contributor 1
  - Contributor 2
* Change 1
  * Change 2
-- Mon, 22 Mar 2010 00:37:30 +0100
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err == nil {
		t.Errorf("should err!=nil, got err=%q\n", err)
	}
}

func TestVariousDateFormat1(t *testing.T) {
	content := `
0.0.1
-- Mon, 22 Mar 2010 00:37:30
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}

	v := s.Versions[0]
	e := "Mon, 22 Mar 2010 00:37:30 +0000"
	if v.Date.Format(DateLayouts[0]) != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.Date.Format(DateLayouts[0]))
	}
	e = "Mon, 22 Mar 2010 00:37:30"
	if v.GetDate() != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.GetDate())
	}
}

func TestVariousDateFormat2(t *testing.T) {
	content := `
0.0.1
-- Mon, 22 Mar 2010
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}

	v := s.Versions[0]
	e := "Mon, 22 Mar 2010 00:00:00 +0000"
	if v.Date.Format(DateLayouts[0]) != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.Date.Format(DateLayouts[0]))
	}
	e = "Mon, 22 Mar 2010"
	if v.GetDate() != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.GetDate())
	}
}

func TestVariousDateFormat3(t *testing.T) {
	content := `
0.0.1
-- Mon 22 Mar 2010
`
	s := Changelog{}
	err := s.Parse([]byte(content))

	if err != nil {
		t.Errorf("should err==nil, got err=%q\n", err)
	}

	v := s.Versions[0]
	e := "Mon, 22 Mar 2010 00:00:00 +0000"
	if v.Date.Format(DateLayouts[0]) != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.Date.Format(DateLayouts[0]))
	}
	e = "Mon 22 Mar 2010"
	if v.GetDate() != e {
		t.Errorf("should s.Date='%q', got s.Date=%q\n", e, v.GetDate())
	}
}
