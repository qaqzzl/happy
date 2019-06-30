package middleware

type Middleware struct {
	Handler		[]func()
}

//type middlewareHandler struct {
//	Handler		func()
//	//async		bool
//}
