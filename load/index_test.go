package load

import (
	"testing"
)

func TestEmptyYamlFile(t *testing.T) {
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

func TestYamlFile1(t *testing.T) {
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

func TestOneVersion(t *testing.T) {
  content := `
name: my package
versions:
  - version: 1.0.0
    date: Mon Jun 27 2016
    changes:
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
  if len(s.Versions[0].Changes)!=3 {
    t.Errorf("should len(s.Versions[0].Changes)='1', got len(s.Versions[0].Changes)=%q\n", len(s.Versions[0].Changes))
  }
  if s.Versions[0].Changes[0]!="change 1" {
    t.Errorf("should s.Versions[0].Changes[0]='change 1', got s.Versions[0].Changes[0]=%q\n", s.Versions[0].Changes[0])
  }
  if s.Versions[0].Changes[1]!="change 2" {
    t.Errorf("should s.Versions[0].Changes[1]='change 2', got s.Versions[0].Changes[1]=%q\n", s.Versions[0].Changes[1])
  }
  if s.Versions[0].Changes[2]!="change 3\nwith multiple lines" {
    t.Errorf("should s.Versions[0].Changes[2]='change 3\nwith multiple lines', got s.Versions[0].Changes[2]=%q\n", s.Versions[0].Changes[2])
  }
}

func TestMissingDate(t *testing.T) {
  content := `
name: my package
versions:
  - version: 1.0.0
    changes:
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
  - version: 0.1.0
  - version: 1.0.0
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
  - version: 0.1.0
  - version: 1.0.0
  - version: 2.0.0
  - name: noversion
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
