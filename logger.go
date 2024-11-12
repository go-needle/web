package web

import (
	"log"
	"time"
)

// Logger is a middleware which defines to log every http request
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s %s in %v", c.StatusCode, c.Method, c.Req.RequestURI, time.Since(t))
	}
}
