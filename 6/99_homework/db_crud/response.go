package main

import (
	"encoding/json"
	"net/http"
)

type MyResponse struct {
	Body map[string]interface{}
}

func NewResponse() *MyResponse {
	body := make(map[string]interface{})
	return &MyResponse{
		Body: body,
	}
}

func (mrsp *MyResponse) ServeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	mrsp.Body["error"] = message
	body, _ := json.Marshal(mrsp.Body)
	w.Write(body)
}

func (mrsp *MyResponse) ServeSuccess(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	body, _ := json.Marshal(mrsp.Body)
	w.Write(body)
}
