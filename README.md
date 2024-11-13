<!-- markdownlint-disable MD033 MD041 -->
<div align="center">

# 🪡web

<!-- prettier-ignore-start -->
<!-- markdownlint-disable-next-line MD036 -->
a web lightweight framework for golang
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
func middleware1() web.HandlerFunc {
	return func(ctx *web.Context) {
		fmt.Println("test1")
		ctx.Next()
		fmt.Println("test4")
	}
}

// Define middleware
func middleware2(ctx *web.Context) {
	fmt.Println("test2")
	ctx.Next()
	fmt.Println("test3")
}

// You need to implement the web.Listener interface
type hello struct{ web.POST } // In this way, you will not need to implement the 'Method()'

func (h *hello) Pattern() string { return "/hello1" }
func (h *hello) Handle() web.HandlerFunc {
	return func(ctx *web.Context) {
		num := ctx.FormData("num")
		fmt.Println(num)
		ctx.JSON(200, web.H{"msg": "hello1"})
	}
}

type response struct {
	Msg string `json:"msg"`
}

func main() {
	// new a client of http
	s := web.Default()
	{
		// define the group and use middleware
		g1 := s.Group("m1").Use(middleware1())
		{
			g2 := g1.Group("m2").Use(middleware2)
			// bind the listener to work pattern in router
			g2.Bind(&hello{}) //  work at POST /m1/m2/hello1
			// also could use this way to add to router
			g2.GET("/hello2", func(ctx *web.Context) {
				fmt.Println(ctx.Query("num"))
				ctx.JSON(200, &response{Msg: "hello2"})
			}) // work at GET /m1/m2/hello2
		}
	}
	// listen on the port
	s.Run(9999)
}
```

## reference
- [gee-web](https://github.com/geektutu/7days-golang/tree/master/gee-web)
- [gin](https://github.com/gin-gonic/gin) 
