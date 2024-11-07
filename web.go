package web

import (
	"fmt"
	"net"
	"net/http"
)

// HandlerFunc defines the request handler used by web
type HandlerFunc func(*Context)

// App implement the interface of ServeHTTP
type App struct {
	router *router
}

func newApp() *App {
	fmt.Println("🪡Welcome to use go-needle-web\n🪡Github: https://github.com/go-needle/web")
	return &App{router: newRouter()}
}

// New is the constructor of web.App
func New() *App {
	return newApp()
}

type Engine struct {
	app *App
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

func getLocalIP() (string, error) {
	adders, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range adders {
		// 检查ip net.IPNet类型
		if inet, ok := address.(*net.IPNet); ok && !inet.IP.IsLoopback() {
			if inet.IP.To4() != nil {
				return inet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("无可用IP地址")
}

// Run defines the method to start a http server
func (app *App) Run(addr string) {
	go fmt.Println("🪡 The web server is running at " + addr)
	err := http.ListenAndServe(addr, &Engine{app})
	if err != nil {
		panic(err)
	}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.app.router.handle(c)
}
