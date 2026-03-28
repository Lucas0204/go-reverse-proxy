package handler

import (
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

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello, world!"))
	})

	return RequestHandler{
		router,
	}
}
