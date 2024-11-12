package web

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type H map[string]any

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	params map[string]string
	// response info
	StatusCode int
	// extra info
	extras map[string]any
	// middlewares and route
	handlers []HandlerFunc
	index    int
	// server pointer
	server *Server
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
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
		c.handlers[c.index](c)
	}
}

// Abort is used in middleware, it means stopping the current middleware
func (c *Context) Abort() {
	c.index = len(c.handlers)
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

func (c *Context) FormData(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) BindJson(obj any) (any, error) {
	err := decodeJSON(c.Req.Body, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (c *Context) Binary() ([]byte, error) {
	var content []byte
	var tmp = make([]byte, 128)
	for {
		n, err := c.Req.Body.Read(tmp)
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
	return c.Req.FormFile(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	if _, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...))); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	if _, err := c.Writer.Write(data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.server.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
