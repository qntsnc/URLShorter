package handlers

import (
	"linkShorter/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RedirectHandler struct {
	Storage storage.Storage
}

func (rh *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortUrl")
	if shortURL == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	url, err := rh.Storage.GetUrl(r.Context(), shortURL)
	if err != nil {
		http.Error(w, "Failed to get URL", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)

}
