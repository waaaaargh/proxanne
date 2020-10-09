package proxanne

import (
	"net/http"
	"net/http/httputil"
	"regexp"
)

type Route struct {
	Matches *regexp.Regexp
	Target  *httputil.ReverseProxy
}

// Router is just a list of Routes
type Router []*Route

// NewRouter initializes a new Router
func NewRouter() Router {
	return make(Router, 0)
}

//
func (r Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	found := false

	for _, route := range r {
		if route.Matches.Match([]byte(req.URL.String())) {
			route.Target.ServeHTTP(res, req)
			return
		}
	}

	if !found {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
