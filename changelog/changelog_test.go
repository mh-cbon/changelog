package changelog

import (
	"testing"
)

func TestParseEmptyYamlFile(t *testing.T) {
	s := Changelog{}
	err := s.Parse([]byte(""))

  if err!=nil {
    t.Errorf("should err=nil, got err=%q\n", err)
  }
  if s.Name!="" {
    t.Errorf("should s.Name='', got s.Name=%q\n", s.Name)
  }
  if s.Author!="" {
    t.Errorf("should s.Author='', got s.Author=%q\n", s.Author)
  }
  if s.Email!="" {
    t.Errorf("should s.Email='', got s.Email=%q\n", s.Email)
  }
}

func TestParseYamlFile1(t *testing.T) {
  content := `
name: name
author: author
email: email (at) some.com
versions:
`
	s := Changelog{}
	err := s.Parse([]byte(content))

  if err!=nil {
    t.Errorf("should err=nil, got err=%q\n", err)
  }
  if s.Name!="name" {
    t.Errorf("should s.Name='name', got s.Name=%q\n", s.Name)
  }
  if s.Author!="author" {
    t.Errorf("should s.Author='author', got s.Author=%q\n", s.Author)
  }
  if s.Email!="email (at) some.com" {
    t.Errorf("should s.Email='email (at) some.com', got s.Email=%q\n", s.Email)
  }
  if len(s.Versions)>0 {
    t.Errorf("should len(s.Versions)='0', got len(s.Versions)=%q\n", len(s.Versions))
  }
}

func TestParseOneVersion(t *testing.T) {
  content := `
name: my package
versions:
  - a_version: 1.0.0
    date: Mon Jun 27 2016
    xupdates:
      - change 1
      - change 2
      - |
        change 3
        with multiple lines`
	s := Changelog{}
	err := s.Parse([]byte(content))

  if err!=nil {
    t.Errorf("should err=nil, got err=%q\n", err)
  }
  if len(s.Versions)!=1 {
    t.Errorf("should len(s.Versions)='1', got len(s.Versions)=%q\n", len(s.Versions))
  }
  if s.Versions[0].Version.String()!="1.0.0" {
    t.Errorf("should s.Versions[0].Version=1.0.0, got s.Versions[0].Version=%q\n", s.Versions[0].Version)
  }
  if s.Versions[0].Date.String()!="Mon Jun 27 2016" {
    t.Errorf("should s.Versions[0].Date='Mon Jun 27 2016', got s.Versions[0].Date=%q\n", s.Versions[0].Date.String())
  }
  updates := s.Versions[0].Updates
  iexpected := 3
  igot := len(updates)
  if iexpected!=igot {
    t.Errorf("should len(s.Versions[0].Updates)=%q, got len(s.Versions[0].Updates)=%q\n", iexpected, igot)
  }
  expected := "change 1"
  got := updates[0]
  if expected!=got {
    t.Errorf("should s.Versions[0].Updates[0]=%q, got s.Versions[0].Updates[0]=%q\n", expected, got)
  }
  expected = "change 2"
  got = updates[1]
  if expected!=got {
    t.Errorf("should s.Versions[0].Updates[1]=%q, got s.Versions[0].Updates[1]=%q\n", expected, got)
  }
  expected = "change 3\nwith multiple lines"
  got = updates[2]
  if expected!=got {
    t.Errorf("should s.Versions[0].Updates[2]=%q, got s.Versions[0].Updates[2]=%q\n", expected, got)
  }
}

func TestParseMissingDate(t *testing.T) {
  content := `
name: my package
versions:
  - a_version: 1.0.0
    xupdates:
      - change 1
`
	s := Changelog{}
	err := s.Parse([]byte(content))

  if err!=nil {
    t.Errorf("should err=nil, got err=%q\n", err)
  }
  if s.Versions[0].Date!=nil {
    t.Errorf("should s.Versions[0].Date=nil, got s.Versions[0].Date=%q\n", s.Versions[0].Date)
  }
}

func TestVersionSort(t *testing.T) {
  content := `
name: my package
versions:
  - a_version: 0.1.0
  - a_version: 1.0.0
`
	s := Changelog{}
	err := s.Parse([]byte(content))

  if err!=nil {
    t.Errorf("should err=nil, got err=%q\n", err)
  }
  s.Sort()
  if s.Versions[0].Version.String()!="1.0.0" {
    t.Errorf("should s.Versions[0].Version=1.0.0, got s.Versions[0].Version=%q\n", s.Versions[0].Version)
  }
  if s.Versions[1].Version.String()!="0.1.0" {
    t.Errorf("should s.Versions[0].Version=0.1.0, got s.Versions[1].Version=%q\n", s.Versions[1].Version)
  }
}

func TestVersionSortExtended(t *testing.T) {
  content := `
name: my package
versions:
  - a_version: 0.1.0
  - a_version: 1.0.0
  - a_version: 2.0.0
  - aname: noversion
`
	s := Changelog{}
	err := s.Parse([]byte(content))

  if err!=nil {
    t.Errorf("should err=nil, got err=%q\n", err)
  }
  s.Sort()
  if s.Versions[0].Version!=nil {
    t.Errorf("should s.Versions[0].Version=nil, got s.Versions[0].Version=%q\n", s.Versions[0].Version)
  }
  if s.Versions[1].Version.String()!="2.0.0" {
    t.Errorf("should s.Versions[1].Version=2.0.0, got s.Versions[1].Version=%q\n", s.Versions[1].Version)
  }
  if s.Versions[2].Version.String()!="1.0.0" {
    t.Errorf("should s.Versions[2].Version=1.0.0, got s.Versions[2].Version=%q\n", s.Versions[2].Version)
  }
  if s.Versions[3].Version.String()!="0.1.0" {
    t.Errorf("should s.Versions[3].Version=0.1.0, got s.Versions[3].Version=%q\n", s.Versions[3].Version)
  }
}

func TestEncode1(t *testing.T) {
  expected := `author: author
email: email
name: name
versions:
- a_version: 0.0.2
  aname: name2
  date: Tue Jun 28 2016
- a_version: 0.0.1
  aname: name
  date: Sat Jun 25 2016
`
	s := Changelog{}
  s.Name = "name"
  s.Author = "author"
  s.Email = "email"
  v1 := s.CreateVersion("name2", "0.0.2", "Tue Jun 28 2016")
  s.Versions = append(s.Versions, v1)
  v2 := s.CreateVersion("name", "0.0.1", "Sat Jun 25 2016")
  s.Versions = append(s.Versions, v2)

  d, err := s.Encode()
  if err!=nil {
    t.Errorf("should err=nil, got err=%q\n", err)
  }
  if string(d)!=expected {
    t.Errorf("should string(d)=%q, got string(d)=%q\n", expected, string(d))
  }
}
