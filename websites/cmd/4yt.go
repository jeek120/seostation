package main

import (
	"github.com/jeek120/seostation/websites/config"
	"github.com/jeek120/seostation/websites/server"
)

func main() {

	c := config.Get()
	s := server.NewServer(c.Listen, server.ServerWithTlsListen(c.TlsListen))

	for _, site := range c.Sites {
		seo := server.NewSeoServer()
		for _, p := range site.ProxyHost {
			seo.WithHost(p.Schema, p.Before, p.After)
		}

		for _, r := range site.Replaces {
			seo.WithReplace(r.Key.(string), r.Value.(string), -1)
		}
		s.AddSite(seo)
	}

	s.Start()
}
