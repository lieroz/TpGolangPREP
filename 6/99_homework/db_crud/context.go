package main

import (
	"net/http"
)

type MyContext struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	PathVars map[string]string
}

func CreateContext(w http.ResponseWriter, r *http.Request) *MyContext {
	return &MyContext{
		Writer:   w,
		Request:  r,
		PathVars: make(map[string]string),
	}
}

func (mctx *MyContext) MatchAndFill(regPath, path []string) bool {
	lenRegPath, lenPath := len(regPath), len(path)
	if lenPath > lenRegPath && len(path[lenPath-1]) == 0 {
		path = path[:lenPath-1]
		lenPath = len(path)
	}
	if lenRegPath != lenPath {
		return false
	}
	for i := 0; i < lenPath; i++ {
		if len(path[i]) != 0 && len(regPath[i]) == 0 {
			return false
		}
		if len(regPath[i]) == 0 || regPath[i][0] != '$' {
			continue
		}
		mctx.PathVars[regPath[i][1:]] = path[i]
	}
	return true
}

func (mctx *MyContext) GetPathVar(pathVar string) string {
	v, ok := mctx.PathVars[pathVar]
	if !ok {
		return ""
	}
	return v
}
