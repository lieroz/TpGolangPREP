package main

import (
	"bytes"
	"testing"
)

const testFullResult = `├───main.go (vary)
├───main_test.go (1318b)
└───testdata
	├───project
	│	├───file.txt (19b)
	│	└───gopher.png (70372b)
	├───static
	│	├───css
	│	│	└───body.css (28b)
	│	├───html
	│	│	└───index.html (57b)
	│	└───js
	│		└───site.js (10b)
	├───zline
	│	└───empty.txt (empty)
	└───zzfile.txt (empty)
`

func TestTreeFull(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, ".", true)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != testFullResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testFullResult)
	}
}

const testDirResult = `└───testdata
	├───project
	├───static
	│	├───css
	│	├───html
	│	└───js
	└───zline
`

func TestTreeDir(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, ".", false)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != testDirResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testDirResult)
	}
}
