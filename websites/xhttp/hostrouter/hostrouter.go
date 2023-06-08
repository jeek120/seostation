package hostrouter

import (
	"net/http"
	"strings"

	"github.com/jeek120/seostation/websites/xhttp"
)

type Routes map[string]xhttp.Router

var _ xhttp.Routes = Routes{}

func New() Routes {
	return Routes{}
}

func (hr Routes) Match(rctx *xhttp.Context, method, path string) bool {
	return true
}

func (hr Routes) Map(host string, h xhttp.Router) {
	hr[strings.ToLower(host)] = h
}

func (hr Routes) Unmap(host string) {
	delete(hr, strings.ToLower(host))
}

func (hr Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := strings.ToLower(requestHost(r))
	if router, ok := hr[host]; ok {
		router.ServeHTTP(w, r)
		return
	}
	if router, ok := hr[getWildcardHost(host)]; ok {
		router.ServeHTTP(w, r)
		return
	}
	if router, ok := hr["*"]; ok {
		router.ServeHTTP(w, r)
		return
	}
	http.Error(w, http.StatusText(404), 404)
}

func (hr Routes) Routes() []xhttp.Route {
	return hr[""].Routes()
}

func (hr Routes) Middlewares() xhttp.Middlewares {
	return xhttp.Middlewares{}
}

func requestHost(r *http.Request) (host string) {
	// not standard, but most popular
	host = r.Header.Get("X-Forwarded-Host")
	if host != "" {
		return
	}

	// RFC 7239
	host = r.Header.Get("Forwarded")
	_, _, host = parseForwarded(host)
	if host != "" {
		return
	}

	// if all else fails fall back to request host
	host = r.Host
	return
}

func parseForwarded(forwarded string) (addr, proto, host string) {
	if forwarded == "" {
		return
	}
	for _, forwardedPair := range strings.Split(forwarded, ";") {
		if tv := strings.SplitN(forwardedPair, "=", 2); len(tv) == 2 {
			token, value := tv[0], tv[1]
			token = strings.TrimSpace(token)
			value = strings.TrimSpace(strings.Trim(value, `"`))
			switch strings.ToLower(token) {
			case "for":
				addr = value
			case "proto":
				proto = value
			case "host":
				host = value
			}

		}
	}
	return
}

func getWildcardHost(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) > 1 {
		wildcard := append([]string{"*"}, parts[1:]...)
		return strings.Join(wildcard, ".")
	}
	return strings.Join(parts, ".")
}
