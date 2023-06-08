package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/jeek120/seostation/websites/xhttp"
	"github.com/jeek120/seostation/websites/xhttp/hostrouter"
	"github.com/jeek120/seostation/websites/xhttp/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type SeoServer struct {
	seos      []*seoServer
	hr        hostrouter.Routes
	listen    string
	tlsListen string
}

type Server interface {
	AddSite(s *seoServer)
	Start()
}

type ServerOpt func(*SeoServer)

func NewServer(listen string, opts ...ServerOpt) Server {
	s := &SeoServer{
		hr:     hostrouter.New(),
		seos:   make([]*seoServer, 0),
		listen: listen,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func ServerWithTlsListen(tlsListen string) ServerOpt {
	return func(s *SeoServer) {
		s.tlsListen = tlsListen
	}
}

func (server *SeoServer) AddSite(s *seoServer) {
	for o, _ := range s.proxyHost {
		server.hr.Map(o, s.seoRouter())
		if !strings.HasPrefix(o, "*") {
			server.hr.Map("*."+o, s.seoRouter())
		}
	}
	server.seos = append(server.seos, s)
}

func (server *SeoServer) Start() {
	r := xhttp.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Requests to api.domain.com

	/*
		s := &seoServer{
			proxyHost: map[string]string{
				"local1.com:8080": "4yt.net",
			},
			replace: []replace{
				{o: "4yt.net/", n: "local1.com:8080/", t: -1},
			},
			proxySchema: "https",
		}
			server.hr.Map("local1.com:8080", s.seoRouter())
			server.hr.Map("*.local1.com:8080", s.seoRouter())
	*/

	r.Mount("/", server.hr)

	if server.tlsListen != "" {
		domains := make([]string, 0)

		for _, seo := range server.seos {
			for domain, _ := range seo.proxyHost {
				domains = append(domains, domain)
			}
		}

		log.Printf("开始申请%v域名的证书", domains)
		http.Serve(autocert.NewListener(domains...), r)
	}

	err := http.ListenAndServe(server.listen, r)
	if err != nil {
		panic(err)
	}
}
