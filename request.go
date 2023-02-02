package clover

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Request http server 请求信息
type Request interface {
	// 原始的http请求
	HTTPRequest() *http.Request

	// Path 获取请求的路径
	// ？前面的那部分
	Path() string

	//Request URI
	RequestURI() string

	// 客户端的地址，如 127.0.0.1:12345
	RemoteAddr() string

	// 请求头信息
	Header(name string) (value string, has bool)

	// 读取header字段
	HeaderDefault(name string, defaultValue string) string

	// 获取原始的cookie
	Cookie(name string) (value *http.Cookie, has bool)

	// Query 获取get请求里的参数
	// 如 xxx?a=v1&b=v2，可获取a、b的值
	Query(name string) (value string, has bool)

	// QueryDefault 获取get请求里的参数，若name不存在，或者是为空字符串，会返回默认值
	QueryDefault(name string, defaultValue string) string

	PostForm(name string) (value string, has bool)

	PostFormDefault(name string, defaultValue string) string

	JsonUnmarshal(dst interface{}) error

	Body() io.ReadCloser
}

// NewRequest 基于原生的request创建一个封装更多功能的request
func NewRequest(req *http.Request) Request {
	return &request{raw: req}
}

type request struct {
	raw      *http.Request
	urlQuery url.Values
}

func (req *request) HTTPRequest() *http.Request {
	return req.raw
}

func (req *request) Path() string {
	return req.raw.URL.Path
}

func (req *request) RequestURI() string {
	return req.raw.RequestURI
}

func (req *request) RemoteAddr() string {
	return req.raw.RemoteAddr
}

func (req *request) Header(name string) (value string, has bool) {
	vs := req.raw.Header.Values(name)
	if len(vs) == 0 {
		return "", false
	}
	return vs[0], true
}

func (req *request) HeaderDefault(name string, defaultValue string) string {
	if v, _ := req.Header(name); v != "" {
		return v
	}
	return defaultValue
}

func (req *request) Cookie(name string) (value *http.Cookie, has bool) {
	cookie, err := req.raw.Cookie(name)
	if err != nil {
		return nil, false
	}
	return cookie, true
}
func (req *request) Query(name string) (value string, has bool) {
	if req.urlQuery == nil {
		req.urlQuery = req.raw.URL.Query()
	}
	values := req.urlQuery[name]
	if len(values) == 0 {
		return "", false
	}
	return values[0], true
}

func (req *request) QueryDefault(name string, defaultValue string) string {
	if v, _ := req.Query(name); v != "" {
		return v
	}
	return defaultValue
}

func (req *request) PostForm(name string) (value string, has bool) {
	_ = req.raw.ParseForm()
	vs := req.raw.PostForm[name]
	if len(vs) == 0 {
		return "", false
	}
	return vs[0], true
}

func (req *request) PostFormDefault(name string, defaultValue string) string {
	if v, _ := req.PostForm(name); v != "" {
		return v
	}
	return defaultValue
}

func (req *request) JsonUnmarshal(dst interface{}) (err error) {
	var bs []byte
	if bs, err = io.ReadAll(req.Body()); err != nil {
		return
	}
	err = json.Unmarshal(bs, dst)
	return
}

func (req *request) Body() io.ReadCloser {
	return req.raw.Body
}

func (req *request) WithContext(ctx context.Context) Request {
	r2 := new(request)
	*r2 = *req
	raw2 := req.raw.WithContext(ctx)
	r2.raw = raw2
	return r2
}

var _ Request = (*request)(nil)
