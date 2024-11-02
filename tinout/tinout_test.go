package tinout

import (
	"os"
	"testing"

	"github.com/rwxrob/bonzai/json"
)

func TestLoad(t *testing.T) {
	spec, err := Load("testdata/commonmark-0.29.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(spec.Name, json.This{spec.Tests})
}

func TestRead(t *testing.T) {
	f, _ := os.Open("testdata/commonmark-0.29.yaml")
	spec, _ := Read(f)
	t.Log(spec.Name, json.This{spec.Tests})
}

func TestCheck(t *testing.T) {
	alltrue := func(t *Test) bool {
		return t.I == t.I // forces them all
	}
	somefalse := func(t *Test) bool {
		t.Got = "nothing"
		return false
	}
	spec, err := Load("testdata/commonmark-0.29.yaml")
	if err != nil {
		t.Fatal(err)
	}
	result := spec.Check(alltrue)
	if result != nil {
		t.Fatal(`all true checks had a failure`)
	}
	result = spec.Check(somefalse)
	if result == nil {
		t.Fatal(`some false checks had no failure`)
	}
	//t.Log(result.State())
}
