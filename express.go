package express

import (
	"net/http"
	"strings"

	rtr "github.com/julienschmidt/httprouter"
)

//__ROUTER

type Params interface{}
type NextFunc func() Handle
//TODO:  Change signature of middlewares so user does not have to write nextfync and return for actual Handle
type Handle func(http.ResponseWriter, *http.Request, Params, NextFunc) Handle

//Route describes the requirements of each subroute
type Route struct {
	Path   string
	Method string
	Func   []Handle
}

func iterateMiddleWare(handlers *[]Handle, res http.ResponseWriter, req *http.Request, prm Params) Handle {
	if len(*handlers) == 1 {
		return (*handlers)[0]
	}
	maxIndex := len(*handlers) - 1
	for index, val := range *handlers {
		hasNext := false
		if index != maxIndex {
			hasNext = true
		}
		if !hasNext {
			return val
		}
		// TODO: Prevent Double execution
		if next := val(res, req, prm, func() Handle { return (*handlers)[index+1] }); next == nil {
			return val
		}
	}
	return nil
}

//Mounter is basic router mount struct
//Rather Than Creating mounter by hand use NewRoutes function which will return a pointer to Mounter
type Mounter struct {
	BasePath string
	Routes   []*Route
}

//New creates the actual router instance
//New must be called before mount which will receive the output of New function
//Most of the times mounter function is required to be called once only
func New() *rtr.Router {
	var router = rtr.New()
	return router
}

//NewRoutes creates the mounter which can be Exported from the package
func NewRoutes(base string) (mounter *Mounter) {
	mounter = &Mounter{
		BasePath: base,
	}
	return
}

//Mount mounts the divided routes to main router
func (mounter *Mounter) Mount(router *rtr.Router) {
	mounter.BasePath = strings.TrimSuffix(mounter.BasePath, "/")
	for _, route := range mounter.Routes {
		fnc := route.Func
		path := route.Path
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		path = mounter.BasePath + path
		switch strings.ToLower(route.Method) {
		case "get":
			router.GET(path, func(res http.ResponseWriter, req *http.Request, params rtr.Params) {
				iterateMiddleWare(&fnc,res,req,params)(res,req,params,nil)
			})
		case "post":
			router.POST(path, func(res http.ResponseWriter, req *http.Request, params rtr.Params) {
				iterateMiddleWare(&fnc,res,req,params)(res,req,params,nil)
			})
		case "delete":
			router.DELETE(path, func(res http.ResponseWriter, req *http.Request, params rtr.Params) {
				iterateMiddleWare(&fnc,res,req,params)(res,req,params,nil)
			})
		case "put":
			router.PUT(path, func(res http.ResponseWriter, req *http.Request, params rtr.Params) {
				iterateMiddleWare(&fnc,res,req,params)(res,req,params,nil)
			})
		}
	}
}

//GET request handler
func (mounter *Mounter) GET(path string, Funcs ...Handle) {

	mounter.Routes = append(mounter.Routes, &Route{path, "get", Funcs})
}

//POST request handler
func (mounter *Mounter) POST(path string, Funcs ...Handle) {
	mounter.Routes = append(mounter.Routes, &Route{path, "post", Funcs})
}

//PUT request handler
func (mounter *Mounter) PUT(path string, Funcs ...Handle) {
	mounter.Routes = append(mounter.Routes, &Route{path, "put", Funcs})
}

//DELETE request handler
func (mounter *Mounter) DELETE(path string, Funcs ...Handle) {
	mounter.Routes = append(mounter.Routes, &Route{path, "delete", Funcs})
}
