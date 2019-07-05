//路由中间件
//路由前缀

package happy

import (
	"net/http"
	"strings"
	"sync"
)

const (
	GET 			= 0
	POST 			= 1
	PUT				= 2
	DELETE			= 3
	CONNECTIBNG	= 4
	HEAD			= 5
	OPTIONS		= 6
	PATCH			= 7
	TRACE			= 8
)

type Route interface {
	RouteAny (path string, handlers ...HandlerFunc)
	RouteGET(path string, handler ...HandlerFunc)
	RoutePOST(path string, handler ...HandlerFunc)
	RoutePUT(path string, handler ...HandlerFunc)
	RouteDELETE(path string, handler ...HandlerFunc)
	RouteHEAD(path string, handler ...HandlerFunc)
	RouteOPTIONS(path string, handler ...HandlerFunc)
	RoutePATCH(path string, handler ...HandlerFunc)
	RouteTRACE(path string, handler ...HandlerFunc)
}

type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]map[string]muxEntry
	//es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}

type muxEntry struct {
	middleware		[]HandlerFunc
	h       		HandlerFunc
	pattern 		string
	groupPrefix		string
}

type RouterGroup struct {
	handlerMiddleware	[]HandlerFunc
	groupPrefix			string
	//ServeMux			ServeMux
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux

// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc array.
type HandlersChain []HandlerFunc

type RouterMux struct {
	serviceMap sync.Map   // map[string]*service
	reqLock    sync.Mutex // protects freeReq
	freeReq    *http.Request
	respLock   sync.Mutex // protects freeResp
	freeResp   *http.Response
}
////注册一个任何请求方法的路由
//func RouteAny(path string, handler func(*Context)) {
//	DefaultServeMux.HandleFunc(path, "GET", handler)
//	DefaultServeMux.HandleFunc(path, "POST", handler)
//	DefaultServeMux.HandleFunc(path, "PUT", handler)
//	DefaultServeMux.HandleFunc(path, "DELETE", handler)
//}

func NewRoute() *RouterGroup{
	return &RouterGroup {

	}
}

//注册一个任何请求方法的路由
func (route RouterGroup) RouteAny(path string, handlers ...HandlerFunc) {
	route.handle(path, "GET", handlers)
	route.handle(path, "POST", handlers)
	route.handle(path, "PUT", handlers)
	route.handle(path, "DELETE", handlers)
	route.handle(path, "PATCH", handlers)
	route.handle(path, "HEAD", handlers)
	route.handle(path, "OPTIONS", handlers)
	route.handle(path, "TRACE", handlers)
}

func (route RouterGroup) RouteGET(path string, handlers ...HandlerFunc) {
	route.handle(path, "GET", handlers)
}

func (route RouterGroup) RoutePOST(path string, handlers ...HandlerFunc) {
	route.handle(path, "POST", handlers)
}

func (route RouterGroup) RoutePUT(path string, handlers ...HandlerFunc) {
	route.handle(path, "PUT", handlers)
}

func (route RouterGroup) RouteDELETE(path string, handlers ...HandlerFunc) {
	route.handle(path, "DELETE", handlers)
}

func (route RouterGroup) RoutePATCH(path string, handlers ...HandlerFunc) {
	route.handle(path, "PATCH", handlers)
}

func (route RouterGroup) RouteHEAD(path string, handlers ...HandlerFunc) {
	route.handle(path, "HEAD", handlers)
}

func (route RouterGroup) RouteOPTIONS(path string, handlers ...HandlerFunc) {
	route.handle(path, "OPTIONS", handlers)
}

func (route RouterGroup) RouteTRACE(path string, handlers ...HandlerFunc) {
	route.handle(path, "TRACE", handlers)
}

//路由分组
func (route *RouterGroup) RouterGroup(groupPrefix string) *RouterGroup {
	groupPrefix = strings.Trim(groupPrefix, "/")
	return &RouterGroup{
		handlerMiddleware: route.handlerMiddleware,
		groupPrefix: groupPrefix,
	}
}

//注册路由中间件
func (route *RouterGroup) UseMiddleware(handler HandlerFunc) *RouterGroup {
	route.handlerMiddleware = append(route.handlerMiddleware, handler)
	return route
}

func (route *RouterGroup) Handle(pattern string, method string, handlers ...HandlerFunc) {
	if handlers == nil {
		panic("http: nil handler")
	}

	//route.handle(pattern, method, handlers)
}

//注册路由
func (route *RouterGroup) handle(pattern string, method string, handler HandlersChain) {
	DefaultServeMux.mu.Lock()
	defer DefaultServeMux.mu.Unlock()
	pattern = route.groupPrefix+"/"+pattern
	if pattern != "/" {
		pattern = strings.Trim(pattern, "/")
	}
	if pattern == "" {
		pattern = "/"
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := DefaultServeMux.m[method][pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if DefaultServeMux.m == nil {
		DefaultServeMux.m = make(map[string]map[string]muxEntry)
	}
	if DefaultServeMux.m[method] == nil {
		DefaultServeMux.m[method] = make(map[string]muxEntry)
	}
	for i:=1; i<len(handler); i++ {
		route.handlerMiddleware = append(route.handlerMiddleware, handler[i])
	}
	e := muxEntry{
		middleware: route.handlerMiddleware,
		h: handler[0],
		pattern: pattern,
	}
	DefaultServeMux.m[method][pattern] = e
	if pattern[0] != '/' {
		DefaultServeMux.hosts = true
	}
}

func (mux *RouterMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{
		ResponseWriter: w,
		Request: r,
	}

	path := strings.Trim(r.URL.Path, "/")

	if path == "" {path = "/"}
	if _,ok := DefaultServeMux.m[r.Method][path]; ok {
		context.Handler = DefaultServeMux.m[r.Method][path].h
		if DefaultServeMux.m[r.Method][path].middleware == nil {
			DefaultServeMux.m[r.Method][path].h(context)
		}
		for _, m := range DefaultServeMux.m[r.Method][path].middleware {
			m(context)
		}
		return
	}
	w.WriteHeader(404)
	return
}

