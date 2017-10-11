package main

import (
	"net"
	"log"
	"os"
	"strings"
	"io"
)

const (
	IndexPage = "index.html"
)

type Server struct {
	Port    string
	WebRoot string
	Workers int
}

func NewServer(port, webRoot string, workers int) *Server {
	return &Server{
		Port:    port,
		WebRoot: webRoot,
		Workers: workers,
	}
}

func (s *Server) ListenAndServe() {
	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalln("server start error:", err)
	}
	dispatcher := NewDispatcher(s.Workers)
	dispatcher.Run()
	log.Println("server started on port:", s.Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("accept connection error:", err)
		}
		job := Job{
			Conn:       conn,
			WorkerFunc: s.serve,
		}
		JobQueue <- job
	}
}

func (s *Server) serve(conn net.Conn) {
	defer conn.Close()
	req, resp := new(Request), &Response{
		Code: StatusOk,
		Description: "OK",
	}
	err := req.Parse(conn)
	if err != nil {
		if err == io.EOF {
			resp.WriteCommonHeaders(conn)
			return
		}
		resp.BuildErrResp(err)
		resp.WriteCommonHeaders(conn)
		return
	}
	var isIndex = strings.HasSuffix(req.AbsPath, "/")
	if isIndex {
		req.AbsPath += IndexPage
	}
	f, err := os.Open(s.WebRoot + req.AbsPath)
	defer f.Close()
	if err != nil {
		if isIndex {
			resp.BuildErrResp(ErrForbidden)
		} else {
			resp.BuildErrResp(ErrNotFound)
		}
		resp.WriteCommonHeaders(conn)
		return
	}
	s.serveMethod(req.Method, resp, conn, f)
}

func (s *Server) serveMethod(method string, resp *Response, conn net.Conn, f *os.File) {
	switch method {
	case Get:
		resp.WriteBody(conn, f)
	case Head:
		resp.Write(conn, f)
	}
}
