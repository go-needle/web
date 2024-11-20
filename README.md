<!-- markdownlint-disable MD033 MD041 -->
<div align="center">

# 🪡web

<!-- prettier-ignore-start -->
<!-- markdownlint-disable-next-line MD036 -->
a lightweight web framework for golang
<!-- prettier-ignore-end -->

<img src="https://img.shields.io/badge/golang-1.21+-blue" alt="golang">
</div>

## installing
Select the version to install

`go get github.com/go-needle/web@version`

If you have already get , you may need to update to the latest version

`go get -u github.com/go-needle/web`


## quickly start
```golang
package main

import (
	"fmt"
	"github.com/go-needle/web"
)

// Define middleware
func middleware1() web.Handler {
	return web.HandlerFunc(func(c *web.Context) {
		fmt.Println("test1")
		c.Next()
		fmt.Println("test4")
	})
}

// Define middleware
func middleware2(c *web.Context) {
	fmt.Println("test2")
	c.Next()
	fmt.Println("test3")
}

// You need to implement the web.Listener interface
type hello struct {
	web.POST // In this way, you will not need to implement the 'Method()'
	cnt      int
}

func (h *hello) Pattern() string { return "/hello1" }
func (h *hello) Handle(c *web.Context) {
	num := c.FormData("num")
	fmt.Println(num)
	h.cnt++
	c.JSON(200, web.H{"msg": "hello1", "cnt": h.cnt})
}

type response struct {
	Msg string `json:"msg"`
}

func main() {
	// new a server of http
	s := web.Default()
	{
		// define the group and use middleware
		g1 := s.Group("m1").Use(middleware1())
		{
			g2 := g1.Group("m2").Use(web.HandlerFunc(middleware2))
			// bind the listener to work pattern in router
			g2.Bind(&hello{}) //  work at POST /m1/m2/hello1
			// also could use this way to add to router
			g2.GET("/hello2", web.HandlerFunc(func(c *web.Context) {
				fmt.Println(c.Query("num"))
				c.JSON(200, &response{Msg: "hello2"})
			})) // work at GET /m1/m2/hello2
		}
	}
	// listen on the port
	s.Run(9999)
}

```
