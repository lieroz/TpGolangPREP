package main

import (
	"net"
	"log"
	"os"
	"strings"
	"sync"
	"io"
)

const (
	IndexPage = "index.html"
)

type Server struct {
	Port    string
	WebRoot string
}

func NewServer(port, webRoot string) *Server {
	return &Server{
		Port:    port,
		WebRoot: webRoot,
	}
}

func (s *Server) ListenAndServe() {
	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalln("server start error:", err)
	}
	dispatcher := NewDispatcher()
	dispatcher.Run()
	log.Println("server started on port:", s.Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("accept connection error:", err)
		}
		//go s.serve(conn)
		job := Job{
			conn:       conn,
			workerFunc: s.serve,
		}
		JobQueue <- job
	}
}

var reqPool = sync.Pool{
	New: func() interface{} {
		return new(Request)
	},
}

var respPool = sync.Pool{
	New: func() interface{} {
		return &Response{
			Code:        StatusOk,
			Description: "OK",
		}
	},
}

func (s *Server) serve(conn net.Conn) {
	defer conn.Close()
	req, resp := reqPool.Get().(*Request), respPool.Get().(*Response)
	defer req.Reset()
	defer reqPool.Put(req)
	defer resp.Reset()
	defer respPool.Put(resp)
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
	var isIndex = strings.HasSuffix(*req.AbsPath, "/")
	if isIndex {
		*req.AbsPath += IndexPage
	}
	f, err := os.Open(s.WebRoot + *req.AbsPath)
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
	s.serveMethod(*req.Method, resp, conn, f)
}

func (s *Server) serveMethod(method string, resp *Response, conn net.Conn, f *os.File) {
	switch method {
	case Get:
		resp.WriteBody(conn, f)
	case Head:
		resp.Write(conn, f)
	}
}
