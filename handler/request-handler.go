package handler

import (
	"go-reverse-proxy/proxy"
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

func NewRequestHandler() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(proxy.Middleware)

	// this route definition serves to prevent 404 in requests to the server
	// and does not have implementation because proxy middleware is responsible
	// for handling all requests (doing proxy pass)
	router.Route("/", func(r chi.Router) {})

	return RequestHandler{
		router,
	}
}
