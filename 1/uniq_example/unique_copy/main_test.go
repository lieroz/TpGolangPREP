package main

import (
	"bufio"
	"strings"
	"testing"
)

var dataOk = `value1
value2
value3`

func TestOk(t *testing.T) {
	br := bufio.NewReader(strings.NewReader(dataOk))
	if err := uniq(br); err != nil {
		t.Errorf("test failed: %v", err)
	}
}

var dataFail = `value2
value1`

func TestFail(t *testing.T) {
	br := bufio.NewReader(strings.NewReader(dataFail))
	if err := uniq(br); err == nil {
		t.Errorf("test failed: %v", err)
	}

}
