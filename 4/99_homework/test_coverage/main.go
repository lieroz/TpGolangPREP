package main

import (
	"net/http/httptest"
	"net/http"
	"fmt"
)

func main() {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL: ts.URL,
	}
	req := SearchRequest{
		Limit: 321,
		Offset: 0,
		Query: "L",
		OrderField: "Age",
		OrderBy: 1,
	}
	fmt.Println(sc.FindUsers(req))
}
