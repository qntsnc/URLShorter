package handlers

import (
	"encoding/json"
	"linkShorter/internal/service/parser"
	"linkShorter/internal/storage"
	"log"
	"net/http"
)

type UrlHandler struct {
	Storage storage.Storage
}

func (h *UrlHandler) PostUrl(w http.ResponseWriter, r *http.Request) {
	var err error
	var d struct {
		Url string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d.Url, err = parser.ParseURL(d.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	str, err := h.Storage.SaveUrl(r.Context(), d.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"shortUrl": str,
	})
}

func (h *UrlHandler) GetUrl(w http.ResponseWriter, r *http.Request) {
	var d struct {
		ShortUrl string `json:"shortUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url, err := h.Storage.GetUrl(r.Context(), d.ShortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Printf("Not found url: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url": url,
	})
}
