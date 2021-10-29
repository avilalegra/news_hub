package web

import (
	"avilego.me/recent_news/env"
	"avilego.me/recent_news/news"
	"html/template"
	"log"
	"net/http"
)

type SearchHandler struct {
	Finder news.Finder
}

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	expr := r.URL.Query().Get("keywords")
	previews := h.Finder.FindRelated(expr)
	var files = []string{
		env.ProjDir() + "/templates/find_news.gohtml",
		env.ProjDir() + "/templates/base.gohtml",
	}
	ts, err := template.New("find_news.gohtml").Funcs(template.FuncMap{
		"unsafe": RenderUnsafe,
	}).ParseFiles(files...)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var data = struct {
		Keywords string
		Previews []news.Preview
	}{
		Keywords: expr,
		Previews: previews,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func RenderUnsafe(s string) template.HTML {
	return template.HTML(s)
}
