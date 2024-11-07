package web

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
)

// HandlerFunc defines the request handler used by web
type HandlerFunc func(*Context)

type Server struct {
	router *router
}

func newServer() *Server {
	fmt.Println("🪡 Welcome to use go-needle-web 🪡")
	return &Server{router: newRouter()}
}

// New is the constructor of web.Server
func New() *Server {
	return newServer()
}

// Engine implement the interface of ServeHTTP
type Engine struct {
	server *Server
}

func (server *Server) addRoute(method string, pattern string, handler HandlerFunc) {
	server.router.addRoute(method, pattern, handler)
}

// REQUEST defines your method to request
func (server *Server) REQUEST(method, pattern string, handler HandlerFunc) {
	server.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (server *Server) GET(pattern string, handler HandlerFunc) {
	server.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (server *Server) POST(pattern string, handler HandlerFunc) {
	server.addRoute("POST", pattern, handler)
}

// PUT defines the method to add PUT request
func (server *Server) PUT(pattern string, handler HandlerFunc) {
	server.addRoute("PUT", pattern, handler)
}

// DELETE defines the method to add DELETE request
func (server *Server) DELETE(pattern string, handler HandlerFunc) {
	server.addRoute("DELETE", pattern, handler)
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
func (server *Server) Run(port int) {
	fmt.Println("🪡 Router total: " + strconv.Itoa(server.router.total))
	ip, err := getLocalIP()
	addr := fmt.Sprintf("%d", port)
	if err == nil {
		fmt.Println("🪡 IP: " + ip + ":" + addr)
	}
	fmt.Println("🪡 The web server is listening at port " + addr)
	err = http.ListenAndServe(":"+addr, &Engine{server})
	if err != nil {
		panic(err)
	}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.server.router.handle(c)
}
