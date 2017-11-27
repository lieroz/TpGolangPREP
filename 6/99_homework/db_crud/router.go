package main

import (
	"strings"
)

type MyRouter struct {
	Routes []*Route
}

type Route struct {
	Path   []string
	Method string
	Func   func(*MyContext, *MyResponse)
}

func NewRouter() *MyRouter {
	return &MyRouter{
		Routes: make([]*Route, 0),
	}
}

func (mrtr *MyRouter) RegisterHandler(path string, method string, foo func(*MyContext, *MyResponse)) {
	route := &Route{
		Path:   strings.Split(path[1:], "/"),
		Method: method,
		Func:   foo,
	}
	mrtr.Routes = append(mrtr.Routes, route)
}
