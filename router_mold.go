package gojimold

import (
	"regexp"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type Route struct {
	Method  string
	Path    string
	RegExp  string
	Handler interface{}
}

type Routes map[string]Route

type MiddlewareMold interface {
	MiddlewareFunc() interface{}
}

type RouterMold struct {
	Routes
	SubRoutes   string
	Middlewares []MiddlewareMold
	HandlerFunc func(Route) interface{}
}

func (rm *RouterMold) Path(name string) string {
	r, ok := rm.Routes[name]
	if !ok {
		return ""
	}
	return r.Path
}

func (rm *RouterMold) Route(name string) Route {
	r, ok := rm.Routes[name]
	if !ok {
		return Route{}
	}
	return r
}

func (rm *RouterMold) Generate() *web.Mux {
	var mux *web.Mux
	if rm.SubRoutes == "" {
		mux = goji.DefaultMux
	} else {
		mux := web.New()
		mux.Use(middleware.RequestID)
		mux.Use(middleware.Recoverer)
		mux.Use(middleware.AutomaticOptions)
		goji.Handle(rm.SubRoutes, mux)
	}

	for _, m := range rm.Middlewares {
		mux.Use(m.MiddlewareFunc())
	}

	var handlerFunc func(Route) interface{}
	if rm.HandlerFunc == nil {
		handlerFunc = func(r Route) interface{} {
			return r.Handler
		}
	} else {
		handlerFunc = rm.HandlerFunc
	}

	for _, r := range rm.Routes {
		var pattern interface{}
		if r.RegExp != "" {
			pattern = regexp.MustCompile(r.RegExp)
		} else {
			pattern = r.Path
		}
		switch r.Method {
		case "HEAD":
			mux.Head(pattern, handlerFunc(r))
		case "GET":
			mux.Get(pattern, handlerFunc(r))
		case "POST":
			mux.Post(pattern, handlerFunc(r))
		case "PUT":
			mux.Put(pattern, handlerFunc(r))
		case "PATCH":
			mux.Patch(pattern, handlerFunc(r))
		case "DELETE":
			mux.Delete(pattern, handlerFunc(r))
		}
	}
	return mux
}
