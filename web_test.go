package web

import (
	"fmt"
	"testing"
	"time"
)

// Define middleware
func middleware1() Handler {
	return HandlerFunc(func(ctx *Context) {
		fmt.Println("test1")
		ctx.Next()
		fmt.Println("test4")
	})
}

// Define middleware
func middleware2(ctx *Context) {
	fmt.Println("test2")
	ctx.Next()
	fmt.Println("test3")
}

// You need to implement the web.Handler interface
type helloHandler struct {
	cnt int
}

func (h *helloHandler) Handle(c *Context) {
	num := c.FormData("num")
	fmt.Println(num)
	h.cnt++
	c.JSON(200, H{"msg": "hello1", "cnt": h.cnt})
}

// You need to implement the web.Listener interface
type hello struct{ POST }            // In this way, you will not need to implement the 'Method()'
func (h *hello) Pattern() string     { return "/hello1" }
func (h *hello) GetHandler() Handler { return &helloHandler{} }

type response struct {
	Msg string `json:"msg"`
}

type Payload struct {
	Name string `json:"name"`
	JWTDefaultParams
}

func Test(t *testing.T) {
	k := []byte("123456")
	// new a server of http
	s := Default()
	{
		// define the group and use middleware
		g1 := s.Group("m1").Use(middleware1())
		{
			g2 := g1.Group("m2").Use(HandlerFunc(middleware2))
			// bind the listener to work pattern in router
			g2.Bind(&hello{}) //  work at POST /m1/m2/hello1
			// also could use this way to add to router
			g2.GET("/hello2", HandlerFunc(func(ctx *Context) {
				fmt.Println(ctx.Query("num"))
				ctx.JSON(200, &response{Msg: "hello2"})
			})) // work at GET /m1/m2/hello2
		}
	}
	s.GET("/login", HandlerFunc(func(c *Context) {
		p := &Payload{Name: "admin"}
		p.Exp = time.Now().Add(time.Second * 20).Unix()
		token, err := CreateToken(k, p)
		if err != nil {
			c.Fail(500, err.Error())
		}
		c.String(200, token)
	}))
	g3 := s.Group("api").Use(JwtConfirm(k, "token", &Payload{}))
	g3.GET("/users", HandlerFunc(func(c *Context) {
		c.JSON(200, c.Extra("jwt").(*Payload))
	}))
	// listen on the port
	s.Run(9999)
}
