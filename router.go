package houston


import (
	"net"
	"path"
)


type RouteHandler func(net.Conn)


type Route struct {
	Path     string
	Handler  RouteHandler
}


type RouterOpts struct {
	ErrorHandler  RouteHandler
}


type Router struct {
	Routes        []Route
	ErrorHandler  RouteHandler
}


func NewRouter(config RouterOpts) Router {
	var newRouter Router
	newRouter = Router{ErrorHandler: config.ErrorHandler}
	return newRouter
}


func BlankRouter() Router {
	return Router{
		ErrorHandler: func(c net.Conn) {
			NotFound(c, "Requested resource inaccessible.")
		},
	}
}


// Get the handler for a given route. If no route matches a
// handler, then return the default error handler.
func (r *Router) GetRouteHandler(targetPath string) RouteHandler {
	cleanedPath := path.Clean(targetPath)

	// TODO Make this match any-length whitespace string, not just single space.
	if targetPath == "" || targetPath == " " {
		return r.GetRouteHandler("/")
	}

	for _, elem := range r.Routes {
		if elem.Path == cleanedPath {
			return elem.Handler
		}
	}

	return r.ErrorHandler
}


func (r *Router) AddRoute(targetPath string, handler RouteHandler) {
	newRoute := Route{path.Clean(targetPath), handler}
	r.Routes = append(r.Routes, newRoute)
}
