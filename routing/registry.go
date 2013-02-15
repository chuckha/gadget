package routing

import (
	"fmt"
	"github.com/redneckbeard/gadget/requests"
	"net/http"
)

var routes []*Route

func Routes(rtes ...*Route) {
	for _, r := range rtes {
		routes = append(routes, r.Flatten()...)
	}
}

func SetIndex(controllerName string) *Route {
	route := newRoute(controllerName)
	route.segment = ""
	route.buildPatterns("")
	return route
}

func Resource(controllerName string, nested ...*Route) *Route {
	route := newRoute(controllerName)
	route.subroutes = nested
	route.buildPatterns("")
	return route
}

func Prefixed(prefix string, nested ...*Route) *Route {
	route := newRoute("")
	route.subroutes = nested
	route.buildPatterns(prefix)
	return route
}

func Handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := requests.New(r)
		for _, route := range routes {
			if route.Match(req) {
				status, body := route.Respond(req)
				if status == 301 || status == 302 {
					http.Redirect(w, r, body.(string), status)
				} else {
					w.WriteHeader(status)
					fmt.Fprint(w, body.(string))
				}
				return
			}
		}
		w.WriteHeader(404)
		fmt.Fprint(w, "")
	}
}