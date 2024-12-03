package web

import (
	"time"
)

// Logger is a middleware which defines to log every http request
func Logger() Handler {
	return HandlerFunc(func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		Log.Infof("[%d] %s %s in %v", c.StatusCode, c.Method, c.Request.RequestURI, time.Since(t))
	})
}
