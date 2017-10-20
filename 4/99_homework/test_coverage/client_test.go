package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
)

func TestBadToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: "some_invalid_access_token",
		URL:         ts.URL,
	}
	req := SearchRequest{}
	_, err := sc.FindUsers(req)
	if err == nil {
		t.Errorf("%s", err)
	}
}

func TestInvalidUrl(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         "some_server_url",
	}
	req := SearchRequest{}
	_, err := sc.FindUsers(req)
	if err == nil {
		t.Errorf("%s", err)
	}
}

func TestLimitAndOffsetOnLessThanZero(t *testing.T) {
	requests := []SearchRequest{
		{
			Limit: -1,
		}, {
			Offset: -5,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	for _, req := range requests {
		_, err := sc.FindUsers(req)
		if err == nil {
			t.Errorf("%s", err)
		}
	}
}

func TestMaxLimit(t *testing.T) {
	requests := []SearchRequest{
		{
			Limit: 100500,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	for _, req := range requests {
		_, err := sc.FindUsers(req)
		if err != nil {
			t.Errorf("%s", err)
		}
	}
}

func TestBadResponse(t *testing.T) {
	requests := []SearchRequest{
		{
			OrderField: "LALALALA",
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	for _, req := range requests {
		_, err := sc.FindUsers(req)
		if err == nil {
			t.Errorf("%s", err)
		}
	}
}

func TestOrderBy(t *testing.T) {
	requests := []SearchRequest{
		{
			OrderField: "Id",
			OrderBy: 1,
		}, {
			OrderField: "Id",
			OrderBy: -1,
		}, {
			OrderField: "Age",
			OrderBy: 1,
		}, {
			OrderField: "Age",
			OrderBy: -1,
		}, {
			OrderBy: 1,
		}, {
			OrderBy: -1,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	for _, req := range requests {
		_, err := sc.FindUsers(req)
		if err != nil {
			t.Errorf("%s", err)
		}
	}
}