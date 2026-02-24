package router

import (
	"linkShorter/internal/handlers"
	"linkShorter/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(storage storage.Storage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	Urlhandler := &handlers.UrlHandler{Storage: storage}
	r.Post("/url", Urlhandler.PostUrl)
	r.Get("/url", Urlhandler.GetUrl)

	return r

}
