package handler

import (
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/handler/api"
	"avilego.me/recent_news/handler/web"
	"net/http"
)

func NewServerHttpHandler() http.Handler {
	mux := http.NewServeMux()
	configRoutes(mux)
	return mux
}

func configRoutes(mux *http.ServeMux) {
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/api/search", api.ApiSearchHandler{Finder: factory.Finder()})
	mux.Handle("/news", web.SearchHandler{Finder: factory.Finder()})
}
