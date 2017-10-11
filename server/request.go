package main

import (
	"net"
	tp "net/textproto"
	"strings"
	"log"
	"bufio"
	"net/url"
	"sync"
	"io"
)

var readerPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}

type Request struct {
	Method  *string
	AbsPath *string
}

func (r *Request) Reset() {
	r.Method = nil
	r.AbsPath = nil
}

func (r *Request) Parse(conn net.Conn) error {
	reader := readerPool.Get().(*bufio.Reader)
	reader.Reset(conn)
	reqLine, err := tp.NewReader(reader).ReadLine()
	defer reader.Reset(nil)
	defer readerPool.Put(reader)
	if err != nil {
		if err == io.EOF {
			return err
		}
		log.Fatalln("error reading connection:", err)
	}
	reqParams := strings.Split(reqLine, " ")
	if !checkMethod(reqParams[0]) {
		return ErrMethodNotAllowed
	}
	if !checkUrl(reqParams[1]) {
		return ErrBadRequest
	}
	r.Method = &reqParams[0]
	if strings.Contains(reqParams[1], "?") {
		path := reqParams[1][:strings.Index(reqParams[1], "?")]
		r.AbsPath = &path
	} else {
		r.AbsPath = &reqParams[1]
	}
	if *r.AbsPath, err = url.QueryUnescape(*r.AbsPath); err != nil {
		log.Fatalln("error decoding query:", err)
	}
	return nil
}

func checkMethod(reqMethod string) bool {
	for _, method := range AllowedMethods {
		if reqMethod == method {
			return true
		}
	}
	return false
}

func checkUrl(reqUrl string) bool {
	if strings.Contains(reqUrl, "../") {
		return false
	}
	return true
}
