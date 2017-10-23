package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"time"
	"reflect"
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

func TestRequests(t *testing.T) {
	expected := SearchResponse{
		Users: []User{
			{
				Id: 0,
				Name: "Boyd Wolf",
				Age: 22,
				About: "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
				Gender: "male",
			}, {
				Id: 1,
				Name: "Hilda Mayer",
				Age: 21,
				About: "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
				Gender: "female",
			}, {
				Id: 2,
				Name: "Brooks Aguilar",
				Age: 25,
				About: "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.\n",
				Gender: "male",
			},
		},
		NextPage: false,
	}
	req := SearchRequest{
		Limit: 3,
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	resp, err := sc.FindUsers(req)
	if err != nil {
		t.Errorf("%s", err)
	}
	if !reflect.DeepEqual(*resp, expected) {
		t.Errorf("invalid result!")
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
			OrderBy:    1,
		}, {
			OrderField: "Id",
			OrderBy:    -1,
		}, {
			OrderField: "Age",
			OrderBy:    1,
		}, {
			OrderField: "Age",
			OrderBy:    -1,
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

func TestCrash(t *testing.T) {
	req := SearchRequest{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	_, err := sc.FindUsers(req)
	if err == nil {
		t.Errorf("%s", err)
	}
}

func TestTimeout(t *testing.T) {
	req := SearchRequest{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 3)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	_, err := sc.FindUsers(req)
	if err == nil {
		t.Errorf("%s", err)
	}
}

func TestBadRequestUnpackJSON(t *testing.T) {
	req := SearchRequest{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	_, err := sc.FindUsers(req)
	if err == nil {
		t.Errorf("%s", err)
	}
}

func TestBadRequestUnknownError(t *testing.T) {
	req := SearchRequest{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg": "lol"}`))
	}))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	_, err := sc.FindUsers(req)
	if err == nil {
		t.Errorf("%s", err)
	}
}

func TestOkUnpackJSON(t *testing.T) {
	req := SearchRequest{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"msg": "lol"}`))
	}))
	defer ts.Close()
	sc := &SearchClient{
		AccessToken: CorrectAccessToken,
		URL:         ts.URL,
	}
	_, err := sc.FindUsers(req)
	if err == nil {
		t.Errorf("%s", err)
	}
}