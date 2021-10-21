package handler

import (
	"avilego.me/recent_news/factory"
	"net/http"
)

func NewServerHttpHandler() http.Handler {
	mux := http.NewServeMux()
	configRoutes(mux)
	return mux
}

func configRoutes(mux *http.ServeMux) {
	mux.Handle("/search", SearchHandler{Finder: factory.Finder()})
}
