package web

import (
	"fmt"
	"strings"
)

func RateLimit(rate int) Handler {
	return HandlerFunc(func(c *Context) {
		remoteAddr := c.Request.RemoteAddr
		forwardedFor := c.GetHeader("X-Forwarded-For")
		if forwardedFor != "" {
			remoteAddr = strings.Split(forwardedFor, ", ")[0]
		}
		fmt.Println(remoteAddr)
	})
}
