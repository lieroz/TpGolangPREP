package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
)

func TestDummy(t *testing.T) {
	t.Errorf("TODO")
}

// сюда писать тесты
func TestSearchServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
}
