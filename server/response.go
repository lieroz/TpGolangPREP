package main

import (
	"net"
	"strings"
	"strconv"
	"time"
	"bytes"
	"os"
	"path/filepath"
	"io"
	"sync"
)

const (
	Base = 10

	HttpVersion   = "HTTP/1.1"
	ServerName    = "tp-autumn-2017-highload"
	HttpSeparator = "\r\n"
	WordSeparator = " "
)

var writerPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}

type Response struct {
	Code        int
	Description string
}

func (r *Response) Reset() {
	r.Code = StatusOk
	r.Description = "OK"
}

func (r *Response) WriteBody(conn net.Conn, f *os.File) {
	buf := r.writeFileInfo(f)
	io.Copy(buf, f)
	conn.Write(buf.Bytes())
	buf.Reset()
	writerPool.Put(buf)
}

func (r *Response) Write(conn net.Conn, f *os.File) {
	buf := r.writeFileInfo(f)
	conn.Write(buf.Bytes())
	buf.Reset()
	writerPool.Put(buf)
}

func (r *Response) writeFileInfo(f *os.File) *bytes.Buffer {
	fileInfo, _ := f.Stat()
	var contentHeaders = [][]string{
		{
			"Content-Length:", strconv.FormatInt(fileInfo.Size(), Base),
		}, {
			"Content-Type:", GetContentType(filepath.Ext(fileInfo.Name())[1:]),
		},
	}
	buf := r.writeCommonHeaders()
	for _, line := range contentHeaders {
		buf.WriteString(strings.Join(line, WordSeparator) + HttpSeparator)
	}
	buf.WriteString(HttpSeparator)
	return buf
}

func (r *Response) WriteCommonHeaders(conn net.Conn) {
	buf := r.writeCommonHeaders()
	conn.Write(buf.Bytes())
	buf.Reset()
	writerPool.Put(buf)
}

func (r *Response) writeCommonHeaders() *bytes.Buffer {
	var commonHeaders = [][]string{
		{
			HttpVersion, strconv.FormatInt(int64(r.Code), Base), r.Description,
		}, {
			"Date:", time.Now().String(),
		}, {
			"Server:", ServerName,
		}, {
			"Connection: Keep-Alive",
		},
	}
	buf := writerPool.Get().(*bytes.Buffer)
	for _, line := range commonHeaders {
		buf.WriteString(strings.Join(line, WordSeparator) + HttpSeparator)
	}
	return buf
}

func (r *Response) BuildErrResp(err error) {
	switch err {
	case ErrBadRequest:
		r.Code = StatusBadRequest
	case ErrForbidden:
		r.Code = StatusForbidden
	case ErrNotFound:
		r.Code = StatusNotFound
	case ErrMethodNotAllowed:
		r.Code = StatusMethodNotAllowed
	}
	r.Description = err.Error()
}
