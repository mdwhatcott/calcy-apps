package httprouter

import (
	"fmt"
	"net/http"
	"strings"
)

func New(routes ...Route) (http.Handler, error) {
	this := &router{root: &treeNode{handlers: &methodHandlers{}}}
	for _, route := range routes {
		err := this.root.Add(route)
		if err != nil {
			return nil, err
		}
	}
	return this, nil
}

type router struct {
	root *treeNode
}

func (this *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.root.Resolve(request.Method, request.URL.Path).ServeHTTP(response, request)
}

type treeNode struct {
	pathElement string
	static      []*treeNode
	handlers    *methodHandlers
}

func (this *treeNode) Add(route Route) error {
	if len(route.Path) == 0 {
		if !this.handlers.Add(route.Method, route.Handler) {
			return fmt.Errorf("invalid route method: %s", route.Method)
		}
		return nil
	}

	if route.Path[0] == '/' {
		route.Path = route.Path[1:]
	} else {
		return fmt.Errorf("invalid route path: %s", route.Path)
	}

	slashIndex := strings.Index(route.Path, "/")
	if slashIndex == -1 {
		slashIndex = len(route.Path)
	}
	node := &treeNode{
		pathElement: route.Path[:slashIndex],
		handlers:    &methodHandlers{},
	}
	this.static = append(this.static, node)
	route.Path = route.Path[slashIndex:]
	return node.Add(route)
}
func (this *treeNode) Resolve(method, path string) http.Handler {
	if len(path) == 0 {
		return this.handlers.Resolve(method)
	}
	path = strings.TrimPrefix(path, "/")
	slashIndex := strings.Index(path, "/")
	if slashIndex == -1 {
		slashIndex = len(path)
	}
	pathElement := path[:slashIndex]
	for _, static := range this.static {
		if static.pathElement == pathElement {
			return static.Resolve(method, strings.TrimPrefix(path, pathElement))
		}
	}
	return http.HandlerFunc(notFoundHandler)
}
func notFoundHandler(response http.ResponseWriter, _ *http.Request) {
	http.Error(response, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

type methodHandlers struct {
	get http.Handler
}

func (this *methodHandlers) Add(method string, handler http.Handler) bool {
	switch method {
	case http.MethodGet:
		this.get = handler
	default:
		return false
	}
	return true
}
func (this *methodHandlers) Resolve(method string) http.Handler {
	switch method {
	case http.MethodGet:
		return this.get
	default:
		return http.HandlerFunc(methodNotAllowedHandler)
	}
}
func methodNotAllowedHandler(response http.ResponseWriter, _ *http.Request) {
	http.Error(response, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

func ParseRoute(method string, path string, handler http.Handler) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
