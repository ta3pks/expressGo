package express

import (
	"net/http"
	"strings"

	rtr "github.com/nikosEfthias/httpRouter"
)

//__ROUTER

func iterateMiddleWare(middleWares *[]MiddleWare, handler Handle, res http.ResponseWriter, req *http.Request, prm Params) Handle {
	if len(*middleWares) == 0 {
		return handler
	}
	for _, val := range *middleWares {
		// TODO: Prevent Double execution
		if next := val(res, req, prm); !next {
			return func(http.ResponseWriter, *http.Request, Params) {

			}
		}
	}
	return handler
}

//New creates the actual router instance
//New must be called before mount which will receive the output of New function
//Most of the times mounter function is required to be called once only
func New() *ServerMux {
	return &ServerMux{rtr.New()}
}

//NewRoutes creates the mounter which can be Exported from the package
func NewRoutes(base string) (mounter *RouterMux) {
	mounter = &RouterMux{
		BasePath: base,
	}
	return
}

//Mount mounts the divided routes to main router
func (mounter *RouterMux) Mount(router *ServerMux) {
	mounter.BasePath = strings.TrimSuffix(mounter.BasePath, "/")
	for _, route := range mounter.SubRoutes {
		fnc := route.Handler
		middleWares := route.MiddleWares
		path := route.Path
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		path = mounter.BasePath + path
		switch strings.ToLower(route.Method) {
		case "get":
			router.GET(path, func(res http.ResponseWriter, req *http.Request, params rtr.Params) {

				iterateMiddleWare(&middleWares, fnc, res, req, httpRouterParamsToExpressParams(params))(res, req, httpRouterParamsToExpressParams(params))
			})
		case "post":
			router.POST(path, func(res http.ResponseWriter, req *http.Request, params rtr.Params) {
				iterateMiddleWare(&middleWares, fnc, res, req, httpRouterParamsToExpressParams(params))(res, req, httpRouterParamsToExpressParams(params))
			})
		case "delete":
			router.DELETE(path, func(res http.ResponseWriter, req *http.Request, params rtr.Params) {
				iterateMiddleWare(&middleWares, fnc, res, req, httpRouterParamsToExpressParams(params))(res, req, httpRouterParamsToExpressParams(params))
			})
		case "put":
			router.PUT(path, func(res http.ResponseWriter, req *http.Request, params rtr.Params) {
				iterateMiddleWare(&middleWares, fnc, res, req, httpRouterParamsToExpressParams(params))(res, req, httpRouterParamsToExpressParams(params))
			})
		}
	}
}

//GET request handler
func (mounter *RouterMux) GET(path string, handler Handle, middlewares ...MiddleWare) {

	mounter.SubRoutes = append(mounter.SubRoutes, &route{path, "get", middlewares, handler})
}

//POST request handler
func (mounter *RouterMux) POST(path string, handler Handle, middlewares ...MiddleWare) {
	mounter.SubRoutes = append(mounter.SubRoutes, &route{path, "post", middlewares, handler})
}

//PUT request handler
func (mounter *RouterMux) PUT(path string, handler Handle, middlewares ...MiddleWare) {
	mounter.SubRoutes = append(mounter.SubRoutes, &route{path, "put", middlewares, handler})
}

//DELETE request handler
func (mounter *RouterMux) DELETE(path string, handler Handle, middlewares ...MiddleWare) {
	mounter.SubRoutes = append(mounter.SubRoutes, &route{path, "delete", middlewares, handler})
}
func httpRouterParamsToExpressParams(params rtr.Params) Params {
	expressParams := make(map[string]string)
	for _, param := range params {
		expressParams[param.Key] = param.Value
	}
	return expressParams
}
