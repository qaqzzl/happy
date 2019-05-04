/*
Package happy implements high performance, minimalist Go web framework.
Example:
package main
import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)
  // Handler
func hello(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}
func main() {
	// Echo instance
	h := happy.Make()
	// Middleware
	h.Use(middleware.Logger())
	h.Use(middleware.Recover())
	// Routes
	h.GET("/", hello)
	// Start server
	h.Logger.Fatal(e.Start(":1323"))
}
Learn more at https://echo.labstack.com
*/

package happy

import (
	"crypto/tls"
	"fmt"
	"github.com/qaqzzl/happy/context"
	"log"
	"net"
	"net/http"
	"regexp"
	"sync"
)

type(
	Happy struct {
		Server				*http.Server
		TLSServer			*http.Server
		Listener			net.Listener
		TLSListener     	net.Listener
		DisableHTTP2		bool					//是否开启http2
		router           	*Router					//路由
		routers          	map[string]*Router		//路由分组

		pool             	sync.Pool
	}

	Controller struct {
		Ctx  *context.Context
	}
)

func New() (h *Happy) {
	h = &Happy{
		Server: new(http.Server),
		TLSServer: new(http.Server),
	}
	//实现ServeHTTP方法
	h.Server.Handler = h
	//初始化路由
	h.router = NewRouter(h)
	h.routers = map[string]*Router{}
	return
}

// Start starts an HTTP server.
func (h *Happy) Send(address string) error {
	h.Server.Addr = address

	return h.startServer(h.Server)
}

// Start starts an HTTPS server.
func (h *Happy) startTLS(address string) error {
	s := h.TLSServer
	s.Addr = address
	if !h.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	return h.startServer(h.Server)
}

// StartServer starts a custom http server. - 开启监听并服务 ListenAndServe
func (h *Happy) startServer(s *http.Server) (err error) {
	if s.TLSConfig == nil {
		if h.Listener == nil {
			h.Listener, err = newListen(s.Addr)
			if err != nil {
				log.Println("listen :", err)
				return err
			}
		}
		//if !.HidePort {
		//	h.colorer.Printf("⇨ http server started on %s\n", e.colorer.Green(e.Listener.Addr()))
		//}
		fmt.Println("开启服务",h.Listener)
		return s.Serve(h.Listener)
	}

	if h.TLSListener == nil {
		l, err := newListen(s.Addr)
		if err != nil {
			return err
		}
		h.TLSListener = tls.NewListener(l, s.TLSConfig)
	}
	//if !e.HidePort {
	//	e.colorer.Printf("⇨ https server started on %s\n", e.colorer.Green(e.TLSListener.Addr()))
	//}
	return s.Serve(h.Listener)
}



// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (h *Happy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := h.Server.Handler
	fmt.Println("handler", handler)
	if handler == nil {
		handler = http.DefaultServeMux
	}
	//if r.RequestURI == "*" && r.Method == "OPTIONS" {
	//	handler = globalOptionsHandler{}
	//}
	fmt.Println("ServeHTTP")

	//v1.0 - 匹配路由
	router :=  h.router.routes
	for key,value := range router {
		regexps := regexp.MustCompile("("+key+")");
		matchs := regexps.FindSubmatch([]byte(r.URL.Path))
		if matchs != nil {
			fmt.Println(value.handler)
			value.handler(context.Contexts{})
		}
	}
	//h.
	//
	//// Execute chain
	//if err := handler(h); err != nil {
	//	e.HTTPErrorHandler(err, c)
	//}

}

// Close immediately stops the server.
// It internally calls `http.Server#Close()`.
func (h *Happy) Close() error {
	if err := h.TLSServer.Close(); err != nil {
		return err
	}
	return h.Server.Close()
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func newListen(address string) (*tcpKeepAliveListener, error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	//return srv.Serve(tcpKeepAliveListener{ln.(net.TCPListener)})
	return &tcpKeepAliveListener{ln.(*net.TCPListener)}, nil
}
//*******************************************************************分割线