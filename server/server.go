package server

import (
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/server/handler"
	"net/http"
)

func NewServerHttpHandler() http.Handler {
	mux := http.NewServeMux()
	configRoutes(mux)
	return mux
}

func configRoutes(mux *http.ServeMux) {
	mux.Handle("/search", handler.SearchHandler{Finder: factory.Finder()})
}
