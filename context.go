package happy

import "net/http"

type Context struct {
	ResponseWriter	http.ResponseWriter
	Request			*http.Request
}