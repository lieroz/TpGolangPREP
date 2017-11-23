package main

import "net/http"
import "net/url"
import "strconv"
import "strings"
import "encoding/json"
import "io/ioutil"
import "io"
import "fmt"

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func parseCrutchyBody(body io.ReadCloser) url.Values {
	b, _ := ioutil.ReadAll(body)
	defer body.Close()
	query := string(b)
	v, _ := url.ParseQuery(query)
	return v
}

func (p *ProfileParams) validateAndFillProfileParams(args url.Values) (err error) {
	p.Login = args.Get("login")
	if len(p.Login) == 0 {
		return fmt.Errorf("login must me not empty")
	}
	if len(p.Login) <= -1 {
		return fmt.Errorf("login len must be >= -1")
	}
	return nil
}

func (p *CreateParams) validateAndFillCreateParams(args url.Values) (err error) {
	p.Login = args.Get("login")
	if len(p.Login) == 0 {
		return fmt.Errorf("login must me not empty")
	}
	if len(p.Login) <= 10 {
		return fmt.Errorf("login len must be >= 10")
	}
	p.Name = args.Get("full_name")
	if len(p.Name) <= -1 {
		return fmt.Errorf("name len must be >= -1")
	}
	p.Status = args.Get("status")
	if len(p.Status) == 0 {
		p.Status = "user"
	}
	if !contains(strings.Split("user|moderator|admin", "|"), p.Status) {
		return fmt.Errorf("status must be one of [user, moderator, admin]")
	}
	if len(p.Status) <= -1 {
		return fmt.Errorf("status len must be >= -1")
	}
	p.Age, err = strconv.Atoi(args.Get("age"))
	if err != nil {
		return fmt.Errorf("age must be int")
	}
	if p.Age <= 0 {
		return fmt.Errorf("age must be >= 0")
	}
	if p.Age >= 128 {
		return fmt.Errorf("age must be <= 128")
	}
	return nil
}


func (srv *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["error"] = ""
	
	
	var v url.Values
	switch r.Method {
	case "POST":
		v = parseCrutchyBody(r.Body)
	default:
		v = r.URL.Query()
	}
	var params ProfileParams
	if err := params.validateAndFillProfileParams(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}
	user, err := srv.Profile(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			w.WriteHeader(err.(ApiError).HTTPStatus)
			resp["error"] = err.Error()
		default:
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = "bad user"
		}
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}
	resp["response"] = user
	body, _ := json.Marshal(resp)
	w.Write(body)
}

func (srv *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["error"] = ""
	
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotAcceptable)
		resp["error"] = "bad method"
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}
	
	if r.Header.Get("X-Auth") != "100500" {
		w.WriteHeader(http.StatusForbidden)
		resp["error"] = "unauthorized"
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}
	var v url.Values
	switch r.Method {
	case "POST":
		v = parseCrutchyBody(r.Body)
	default:
		v = r.URL.Query()
	}
	var params CreateParams
	if err := params.validateAndFillCreateParams(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}
	user, err := srv.Create(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			w.WriteHeader(err.(ApiError).HTTPStatus)
			resp["error"] = err.Error()
		default:
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = "bad user"
		}
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}
	resp["response"] = user
	body, _ := json.Marshal(resp)
	w.Write(body)
}

func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	switch r.URL.Path {
	
	case "/user/profile":
		srv.handlerProfile(w, r)

	case "/user/create":
		srv.handlerCreate(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
		resp["error"] = "unknown method"
		body, _ := json.Marshal(resp)
		w.Write(body)
	}
}
