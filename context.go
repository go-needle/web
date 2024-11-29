package web

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
)

type H map[string]any

type Context struct {
	// origin objects
	Writer  http.ResponseWriter
	Request *http.Request
	// request info
	Path   string
	Method string
	params map[string]string
	// response info
	StatusCode int
	// extra info
	extras map[string]any
	// middlewares and route
	handlers []Handler
	index    int
	// server pointer
	server *Server
	// response tag
	isResponse bool
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:    req.URL.Path,
		Method:  req.Method,
		Request: req,
		Writer:  w,
		extras:  make(map[string]any),
		index:   -1,
	}
}

func decodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

// Next is used in middleware, it means executing the next middleware or handle
func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index].Handle(c)
	}
}

// Abort is used in middleware, it means stopping the current middleware
func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// Abort is used in middleware, it means stopping the current middleware
func (c *Context) isAbort() bool {
	return c.index == len(c.handlers)
}

// Param is used to get the parameter at path.
func (c *Context) Param(key string) string {
	value, _ := c.params[key]
	return value
}

// Extra is used to get the info which set by user
func (c *Context) Extra(key string) any {
	value, _ := c.extras[key]
	return value
}

// SetExtra is used to set the info by user
func (c *Context) SetExtra(key string, v any) {
	c.extras[key] = v
}

func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context)  ClientIp() string{
	remoteAddr := c.Request.RemoteAddr
	forwardedFor := c.GetHeader("X-Forwarded-For")
	if forwardedFor != "" {
		remoteAddr = strings.Split(forwardedFor, ",")[0]
	}
	return net.ParseIP(remoteAddr).String()
}

func (c *Context) FormData(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) BindJson(obj any) (any, error) {
	err := decodeJSON(c.Request.Body, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (c *Context) Binary() ([]byte, error) {
	var content []byte
	var tmp = make([]byte, 128)
	for {
		n, err := c.Request.Body.Read(tmp)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		content = append(content, tmp[:n]...)
	}
	return content, nil
}

func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...any) {
	if c.isResponse {
		return
	}
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	if _, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...))); err != nil {
		panic(err)
	}
	c.isResponse = true
}

func (c *Context) JSON(code int, obj any) {
	if c.isResponse {
		return
	}
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
	c.isResponse = true
}

func (c *Context) Data(code int, contentType string, data []byte) {
	if c.isResponse {
		return
	}
	c.SetHeader("Content-Type", contentType)
	c.Status(code)
	if _, err := c.Writer.Write(data); err != nil {
		panic(err)
	}
	c.isResponse = true
}

func (c *Context) HTML(code int, name string, data interface{}) {
	if c.isResponse {
		return
	}
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.server.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		panic(err)
	}
	c.isResponse = true
}

func (c *Context) Fail(code int, err string) {
	c.Abort()
	c.String(code, err)
}
