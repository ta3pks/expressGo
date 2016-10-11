package express

import (
	"net/http"

	"github.com/nikosEfthias/httpRouter"
)

type route struct {
	Path        string
	Method      string
	MiddleWares []MiddleWare
	Handler     Handle
}

//ServerMux is created once and used with listenAndServe function
type ServerMux struct {
	*httprouter.Router
}

//RouterMux is the multiplexer for mainRoutes
type RouterMux struct {
	SubRoutes []*route
	BasePath  string
}

//Params are keyvalue pairs which can be read from the path
//e.g. /test/path/:name //Params{name:"<whatever>"}
type Params map[string]string

//NextFunc calls the next middleware if theres no middleware it finally calls the handler
//if next is not called the middleware has to return -1 and respond to user
//if middleware does not respond to user in that case connecton will simply hang thats what we dont want generally
//if theres a problem with the pre checks we do we can inform user from middleware and stop handling the request
type NextFunc func() bool

//MiddleWare functions can be chained and executed before the request handler is executed
//Middlewares are great to do things like password control etc. instead of writing the same thign all the time
type MiddleWare func(http.ResponseWriter, *http.Request, Params) bool

//Handle Handles the request
type Handle func(http.ResponseWriter, *http.Request, Params)
