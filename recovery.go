package web

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// Recovery is a middleware which defines to prevent panic from causing HTTP service termination
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
				log.Printf("[%d] %s %s Internal Server Error", c.StatusCode, c.Method, c.Req.RequestURI)
				fmt.Printf("\033[31m%s\n\n\033[0m", trace(message))
			}
		}()

		c.Next()
	}
}
