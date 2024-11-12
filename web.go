package web

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"path"
	"strconv"
	"time"
)

// HandlerFunc defines the request handler used by web
type HandlerFunc func(*Context)

type Listener interface {
	Method() string
	Pattern() string
	Handle() HandlerFunc
}

type GET struct{}

func (*GET) Method() string { return "GET" }

type POST struct{}

func (*POST) Method() string { return "POST" }

type DELETE struct{}

func (*DELETE) Method() string { return "DELETE" }

type PUT struct{}

func (*PUT) Method() string { return "PUT" }

type PATCH struct{}

func (*PATCH) Method() string { return "PATCH" }

type OPTIONS struct{}

func (*OPTIONS) Method() string { return "OPTIONS" }

type HEAD struct{}

func (*HEAD) Method() string { return "HEAD" }

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	server      *Server       // all groups share a Server instance
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	if len(prefix) == 1 {
		panic("the length of prefix must > 0")
	}
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	server := group.server
	groupPrefix := group.prefix + prefix
	newGroup := &RouterGroup{
		prefix: groupPrefix,
		parent: group,
		server: server,
	}
	server.groups.insert(groupPrefix, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.server.router.addRoute(method, pattern, handler)
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) *RouterGroup {
	group.middlewares = append(group.middlewares, middlewares...)
	return group
}

// Bind is defined to bind all listeners to the router
func (group *RouterGroup) Bind(listeners ...Listener) {
	for _, listener := range listeners {
		group.REQUEST(listener.Method(), listener.Pattern(), listener.Handle())
	}
}

// REQUEST defines your method to request
func (group *RouterGroup) REQUEST(method, pattern string, handler HandlerFunc) {
	if len(pattern) == 1 {
		panic("the length of pattern must > 0")
	}
	if pattern[0] != '/' {
		pattern = "/" + pattern
	}
	group.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.REQUEST("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.REQUEST("POST", pattern, handler)
}

// PUT defines the method to add PUT request
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.REQUEST("PUT", pattern, handler)
}

// DELETE defines the method to add DELETE request
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.REQUEST("DELETE", pattern, handler)
}

// PATCH defines the method to add PATCH request
func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.REQUEST("PATCH", pattern, handler)
}

// OPTIONS defines the method to add OPTIONS request
func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFunc) {
	group.REQUEST("OPTIONS", pattern, handler)
}

// HEAD defines the method to add HEAD request
func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	group.REQUEST("HEAD", pattern, handler)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Fail(http.StatusNotFound, err.Error())
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static is defined to map local static resources
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

type Server struct {
	*RouterGroup
	router        *router
	groups        *trieTreeG         // store all groups
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

func newServer() *Server {
	server := &Server{router: newRouter()}
	server.RouterGroup = &RouterGroup{server: server}
	server.groups = newTrieTreeG(server.RouterGroup)
	return server
}

// New is the constructor of web.Server
func New() *Server {
	return newServer()
}

// Default is the constructor of web.Server with Recovery and Logger
func Default() *Server {
	server := newServer()
	server.Use(Recovery(), Logger())
	return server
}

func (server *Server) SetFuncMap(funcMap template.FuncMap) {
	server.funcMap = funcMap
}

func (server *Server) LoadHTMLGlob(pattern string) {
	server.htmlTemplates = template.Must(template.New("").Funcs(server.funcMap).ParseGlob(pattern))
}

// Engine implement the interface of ServeHTTP
type Engine struct {
	server *Server
}

func getInternalIP() (string, error) {
	adders, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range adders {
		// 检查IP地址是否是定义在回环接口上的IPv4地址
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no internal IP address found, check for multiple interfaces")
}

func welcome(port string, router int) {
	time.Sleep(time.Millisecond * 100)
	fmt.Println("🪡 Welcome to use go-needle-web")
	fmt.Println("🪡 Available router total: " + strconv.Itoa(router))
	ip, err := getInternalIP()
	if err == nil {
		fmt.Println("🪡 IP address: " + ip)
	}
	fmt.Println("🪡 The web server is listening at port " + port)
}

// Run defines the method to start a http server
func (server *Server) Run(port int) {
	portStr := strconv.Itoa(port)
	go welcome(portStr, server.router.total)
	err := http.ListenAndServe(":"+portStr, &Engine{server})
	if err != nil {
		panic(err)
	}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	middlewaresFind := engine.server.groups.search(req.URL.Path)
	c := newContext(w, req)
	c.handlers = middlewaresFind
	c.server = engine.server
	engine.server.router.handle(c)
}
