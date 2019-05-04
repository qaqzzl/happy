//正则路由
//路由分组
package happy

import (
	"fmt"
	"github.com/qaqzzl/happy/context"
	"net/http"
	"reflect"
	"runtime"
)

type (
	Router struct {
		//tree   *node
		routes map[string]*Route
		//echo   *Echo
	}

	// Route contains a handler and information for matching against requests.
	Route struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Name   string `json:"name"`
		handler	func(ctx context.Contexts) error
	}

	routes struct {
		match string
		methodHandler *methodHandler
		handler func()
	}
)

type methodHandler struct {
	connect  context.HandlerFunc
	delete   context.HandlerFunc
	get      context.HandlerFunc
	head     context.HandlerFunc
	options  context.HandlerFunc
	patch    context.HandlerFunc
	post     context.HandlerFunc
	propfind context.HandlerFunc
	put      context.HandlerFunc
	trace    context.HandlerFunc
}

// NewRouter returns a new Router instance.
func NewRouter(h *Happy) *Router {
	return &Router{
		//tree: &node{
		//	methodHandler: new(methodHandler),
		//},
		routes: map[string]*Route{},
		//echo:   e,
	}
}

//注册一个 GET 路由
func (h *Happy) RouteGet(path string, handler context.HandlerFunc) *Route {
	return h.RouteAdd(http.MethodGet, path, handler)
}

//注册一个 POST 路由

//注册一个 DELETE 路由
func RouteDelete() {

}

//注册一个 PUT 路由
func RoutePut() {

}

//注册一个 Any 路由
func (h *Happy) RouteAny(path string, handler context.HandlerFunc) *Route {
	return h.RouteAdd("", path, handler)
}

//	, middleware ...MiddlewareFunc
func (h *Happy) RouteAdd(method string, path string, handler context.HandlerFunc) *Route {
	return h.routeAdd("", method, path, handler)
}

//  , middleware ...MiddlewareFunc
func (h *Happy) routeAdd(host string, method string, path string, handler context.HandlerFunc) *Route {
	name := handlerName(handler)
	//router := e.findRouter(host)
	//router.Add(method, path, func(c Context) error {
	//	h := handler
	//	// Chain middleware
	//	for i := len(middleware) - 1; i >= 0; i-- {
	//		h = middleware[i](h)
	//	}
	//	return h(c)
	//})
	r := &Route{
		Method: method,
		Path:   path,
		Name:   name,
		handler: handler,
	}
	fmt.Println(h.router)
	h.router.routes[method+path] = r
	return r
}

func handlerName(h context.HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}