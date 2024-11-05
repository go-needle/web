package web

import (
	"net/http"
)

// HandlerFunc defines the request handler used by web
type HandlerFunc func(*Context)

// App implement the interface of ServeHTTP
type App struct {
	router *router
}

type Engine struct {
	app *App
}

// New is the constructor of web.App
func New() *App {
	return &App{router: newRouter()}
}

func (app *App) addRoute(method string, pattern string, handler HandlerFunc) {
	app.router.addRoute(method, pattern, handler)
}

// REQUEST defines your method to request
func (app *App) REQUEST(method, pattern string, handler HandlerFunc) {
	app.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (app *App) GET(pattern string, handler HandlerFunc) {
	app.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (app *App) POST(pattern string, handler HandlerFunc) {
	app.addRoute("POST", pattern, handler)
}

// PUT defines the method to add PUT request
func (app *App) PUT(pattern string, handler HandlerFunc) {
	app.addRoute("PUT", pattern, handler)
}

// DELETE defines the method to add DELETE request
func (app *App) DELETE(pattern string, handler HandlerFunc) {
	app.addRoute("DELETE", pattern, handler)
}

// Run defines the method to start a http server
func (app *App) Run(addr string) {
	panic(http.ListenAndServe(addr, &Engine{app}))
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.app.router.handle(c)
}
