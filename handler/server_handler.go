package handler

import (
	"avilego.me/recent_news/env"
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/handler/api"
	"avilego.me/recent_news/handler/web"
	"html/template"
	"log"
	"net/http"
)

func NewServerHttpHandler() http.Handler {
	mux := http.NewServeMux()
	configRoutes(mux)
	return mux
}

func configRoutes(mux *http.ServeMux) {
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/api/news", api.SearchHandler{Finder: factory.Finder()})
	mux.Handle("/news", web.SearchHandler{Finder: factory.Finder()})
	mux.Handle("/", http.HandlerFunc(rootHandler))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.RedirectHandler("/news", 301)
	} else {
		notFoundHandler(w, r)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	ts, err := template.New("404.gohtml").ParseFiles(
		env.ProjDir()+"/templates/base.gohtml",
		env.ProjDir()+"/templates/404.gohtml",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}

	w.WriteHeader(404)
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
