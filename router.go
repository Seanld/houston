package houston


import (
	"path"
)


// Sandboxes represent "static" file directories
// where having code execute upon URL visitation
// is not necessary.
type Sandbox struct {
	Path string
	LocalPath string
}


type RouteHandler func(Context)


type Route struct {
	Path     string
	Handler  RouteHandler
}


type RouterOpts struct {
	ErrorHandler  RouteHandler
}


type Router struct {
	Routes        []Route
	Sandboxes     []Sandbox
	ErrorHandler  RouteHandler
}


func NewRouter(config RouterOpts) Router {
	var newRouter Router
	newRouter = Router{ErrorHandler: config.ErrorHandler}
	return newRouter
}


func BlankRouter() Router {
	return Router{
		ErrorHandler: func(ctx Context) {
			ctx.NotFound("Requested resource inaccessible.")
		},
	}
}


// Get the handler for a given route. If no route matches a
// handler, then return the default error handler.
func (r *Router) GetRouteHandler(targetPath string) RouteHandler {
	cleanedPath := path.Clean(targetPath)

	if cleanedPath == "" || cleanedPath == " " || cleanedPath == "." {
		return r.GetRouteHandler("/")
	}

	// If no static sandboxed files, then find
	// a route.
	for _, elem := range r.Routes {
		if elem.Path == cleanedPath {
			return elem.Handler
		}
	}

	return r.ErrorHandler
}


func (r *Router) Handle(targetPath string, handler RouteHandler) {
	newRoute := Route{path.Clean(targetPath), handler}
	r.Routes = append(r.Routes, newRoute)
}


func (r *Router) Sandbox(targetPath string, sandboxDirPath string) {
	newSandbox := Sandbox{path.Clean(targetPath), sandboxDirPath}
	r.Sandboxes = append(r.Sandboxes, newSandbox)
}
