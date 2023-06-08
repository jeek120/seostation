package server

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/jeek120/seostation/websites/util"
	"github.com/jeek120/seostation/websites/xhttp"
)

type seoServer struct {
	proxyHost map[string]proxyHost
	replace   []replace
}

type proxyHost struct {
	before string
	after  string
	schema string
}
type replace struct {
	o string
	n string
	t int
}

// Router for the API service
func (s *seoServer) seoRouter() xhttp.Router {
	r := xhttp.NewRouter()
	r.Mount("/", s)
	return r
}

func (s *seoServer) seoHandler(w http.ResponseWriter, r *http.Request) {
	p := s.findNewHost(r.Host)
	proxyUrl := p.schema + "://" + p.after + r.RequestURI
	bs, err := util.HttpGetHeader(proxyUrl, nil)
	if err != nil {
		log.Printf("请求代理地址失败: %s url=%s", err, proxyUrl)
		panic(err)
	}
	for _, r := range s.replace {
		bs = bytes.Replace(bs, []byte(r.o), []byte(r.n), r.t)
	}
	w.Write(bs)
}

func (s *seoServer) findNewHost(host string) proxyHost {
	if h, ok := s.proxyHost[host]; ok {
		return h
	}
	topHost := util.GetTopDomain(host)
	if h, ok := s.proxyHost[topHost]; ok {
		s1 := strings.Replace(host+"/", topHost+"/", h.after+"/", 1)
		return proxyHost{
			before: host,
			after:  s1[0 : len(s1)-1],
			schema: h.schema,
		}
	}
	panic("没有找到 " + host + " 的映射")
}

func (s *seoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.seoHandler(w, r)
}

func NewSeoServer(opts ...OptsFunc) *seoServer {
	s := &seoServer{
		proxyHost: make(map[string]proxyHost, 0),
		replace:   make([]replace, 0),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type OptsFunc func(s *seoServer)

func (s *seoServer) WithHost(schema, o string, n string) *seoServer {
	WithHost(schema, o, n)(s)
	return s
}
func WithHost(schema, o string, n string) OptsFunc {
	return func(s *seoServer) {
		s.proxyHost[o] = proxyHost{
			before: o,
			after:  n,
			schema: schema,
		}
	}
}

func (s *seoServer) WithReplace(o, n string, t int) *seoServer {
	WithReplace(o, n, t)(s)
	return s
}
func WithReplace(o, n string, t int) OptsFunc {
	return func(s *seoServer) {
		s.replace = append(s.replace, replace{o: o, n: n, t: t})
	}
}
