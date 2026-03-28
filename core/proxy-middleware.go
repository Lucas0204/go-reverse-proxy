package core

import (
	configTypes "go-reverse-proxy/config-file"
	"go-reverse-proxy/parser"
	"net/http"
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
	// now I can use config from config file
	// m.config.FindBlocksByName(...)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-From-Proxy", "true")
		block := m.config.FindBlocksByName("http")
		parser.PrintBlock(block[0], 0)
		w.Write([]byte("Hello, world!"))
	})
}
