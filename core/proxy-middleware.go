package core

import (
	configTypes "go-reverse-proxy/config-file"
	"go-reverse-proxy/parser"
	"net/http"
	"net/url"
	"strings"
)

type ProxyMiddleware struct {
	config *configTypes.Config
}

func NewProxyMiddleware(configFilePath string) *ProxyMiddleware {
	config := parser.ParseConfig(configFilePath)
	return &ProxyMiddleware{
		config,
	}
}

func (m *ProxyMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-From-Proxy", "true")
		//block := m.config.FindBlocksByName("http")
		//parser.PrintBlock(block[0], 0)

		// get a location block that matches the request
		locationBlock := m.findMatchLocation(r.URL)
		if locationBlock != nil {
			w.Write([]byte("Proxy passing!"))
			return
		}

		w.Write([]byte("Hello from xnign!"))
	})
}

func (m *ProxyMiddleware) findMatchLocation(requestLocation *url.URL) *configTypes.Block {
	locationBlocks := m.config.FindBlocksByName("location")
	for _, block := range locationBlocks {
		location := block.Args[0]
		if strings.HasPrefix(requestLocation.Path, location) {
			return block
		}
	}

	return nil
}
