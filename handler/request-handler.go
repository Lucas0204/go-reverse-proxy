package handler

import (
	"go-reverse-proxy/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RequestHandler struct {
	router *chi.Mux
}

func (rh RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.router.ServeHTTP(w, r)
}

func NewRequestHandler(proxyMiddleware *core.ProxyMiddleware) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(proxyMiddleware.Middleware)

	// this route definition serves to prevent 404 in requests to the server
	// and does not have implementation because core middleware is responsible
	// for handling all requests (doing core pass)
	router.Route("/", func(r chi.Router) {})

	return RequestHandler{
		router,
	}
}
