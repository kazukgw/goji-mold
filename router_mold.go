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

type Routes []Route

type MiddlewareMold interface {
	MiddlewareFunc() interface{}
}

type RouterMold struct {
	Routes
	SubRoutes   string
	Middlewares []MiddlewareMold
	HandlerFunc func(Route) interface{}
}

func (rg *RouterMold) Generate() *web.Mux {
	var mux *web.Mux
	if rg.SubRoutes == "" {
		mux = goji.DefaultMux
	} else {
		mux := web.New()
		mux.Use(middleware.RequestID)
		mux.Use(middleware.Recoverer)
		mux.Use(middleware.AutomaticOptions)
		goji.Handle(rg.SubRoutes, mux)
	}

	for _, m := range rg.Middlewares {
		mux.Use(m.MiddlewareFunc())
	}

	var handlerFunc func(Route) interface{}
	if rg.HandlerFunc == nil {
		handlerFunc = func(r Route) interface{} {
			return r.Handler
		}
	} else {
		handlerFunc = rg.HandlerFunc
	}

	for _, r := range rg.Routes {
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
