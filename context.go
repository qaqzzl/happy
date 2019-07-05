package happy

import "net/http"

type Context struct {
	ResponseWriter	http.ResponseWriter
	Request			*http.Request
	Handler			HandlerFunc
	Test			interface{}
}

func (c *Context) PostFormValue(val string) {
	c.Request.PostFormValue("act")
}