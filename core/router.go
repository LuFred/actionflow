package core

import (
	"actionflow/pkg/osutil"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type RouteVerbType int8

const (
	Get RouteVerbType = iota - 1
	Post
	Put
	Delete
	Patch
	Options
	Trace
	Head
	Connect
)

type Router struct {
	chi *chi.Mux
	s   *Server
}

func NewRouter(s *Server) *Router {
	r := &Router{
		chi: chi.NewRouter(),
		s:   s,
	}
	return r
}

func (r *Router) UseMiddlewares(middlewares ...func(*Server, http.Handler) http.Handler) {
	for _, v := range middlewares {
		r.chi.Use(addMiddlewares(r.s, v))
	}
}

func addMiddlewares(s *Server, mid func(*Server, http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return mid(s, next)
	}
}

func (r *Router) UseHandlers(rs ...Route) {
	lg := r.s.Logger()
	for _, route := range rs {
		if len(route.Pattern) == 0 || !strings.HasPrefix(route.Pattern, "/") {
			lg.Error(fmt.Sprintf(" routing pattern must begin with '/' in %s", route.Pattern))
			osutil.Exit(1)
		}
		if len(route.Routes) > 0 {
			r.chi.Route(route.Pattern, func(cr chi.Router) {
				for _, rts := range route.Routes {
					switch rts.Verb {
					case Get:
						cr.Get(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					case Post:
						cr.Post(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					case Put:
						cr.Put(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					case Patch:
						cr.Patch(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					case Connect:
						cr.Connect(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					case Delete:
						cr.Delete(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					case Options:
						cr.Options(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					case Head:
						cr.Head(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					case Trace:
						cr.Trace(rts.Path, CallHandler(r.s, rts.Dto, route.Handler, rts.Method))
					}
				}
			})
		}
	}
}

type RouteDescribe struct {
	Verb   RouteVerbType
	Path   string
	Method string
	Dto    interface{}
}

type Route struct {
	Pattern string
	Handler HandlerInterface
	Routes  []RouteDescribe
}
