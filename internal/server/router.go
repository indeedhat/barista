package server

import (
	"net/http"
	"path"
	"strings"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type ServerConfig struct {
	MaxBodySize int64
}

type Router struct {
	mux        *http.ServeMux
	middleware []Middleware
	basePath   string
	cfg        ServerConfig
}

// NewRouter creates a new router instance with the provided middleware stack assigned
func NewRouter(serverCfg ServerConfig, middleware ...Middleware) Router {
	return Router{
		mux:        http.DefaultServeMux,
		middleware: middleware,
		cfg:        serverCfg,
	}
}

// ServerMux returns the underlying http.ServeMux instance
func (r Router) ServerMux() *http.ServeMux {
	return r.mux
}

// All registers a handler on all request methods on the provided uri
func (r Router) HandleFunc(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc(r.prefix(path), r.wrap(handler, middleware...))
}

func (r Router) Handle(path string, handler http.Handler, middleware ...Middleware) {
	r.mux.HandleFunc(r.prefix(path), r.wrap(handler.ServeHTTP, middleware...))
}

// Group creates a sub router and assigns a base path and middleware to all routes assigned within it
func (r Router) Group(path string, middleware ...Middleware) Router {
	return Router{
		mux:        r.mux,
		basePath:   r.basePath + path,
		middleware: append(r.middleware, middleware...),
		cfg:        r.cfg,
	}
}

// apply middleware to the handler
func (r Router) apply(handler http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	stack := append(r.middleware, middleware...)

	for i := range stack {
		if stack[len(stack)-1-i] == nil {
			continue
		}

		handler = stack[len(stack)-1-i](handler)
	}

	return handler
}

// wrap handler in a response parsing closure
func (r Router) wrap(handler http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := NewContext(req.Context(), r.cfg)
		req = req.WithContext(ctx)

		r.apply(handler, middleware...)(rw, req)
	}
}

func (r Router) prefix(uri string) string {
	suffix := ""
	if uri[len(uri)-1] == '/' {
		suffix = "/"
	}
	if strings.Contains(uri, " ") {
		parts := strings.SplitN(uri, " ", 2)
		return parts[0] + " " + path.Join(r.basePath, strings.Trim(parts[1], " ")) + suffix
	}

	return path.Join(r.basePath, uri) + suffix
}
