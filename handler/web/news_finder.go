package web

import (
	"avilego.me/recent_news/env"
	"avilego.me/recent_news/news"
	"html/template"
	"net/http"
)

type SearchHandler struct {
	Finder news.Finder
}

const latestNewsCount = 50

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var previews []news.Preview
	expr := r.URL.Query().Get("keywords")

	if expr == "" {
		previews = h.Finder.FindLatest(latestNewsCount)
	} else {
		previews = h.Finder.FindRelated(expr)
	}

	ts, err := template.New("find_news.gohtml").Funcs(template.FuncMap{
		"unsafe": RenderUnsafe,
	}).ParseFiles(
		env.ProjDir()+"/templates/base.gohtml",
		env.ProjDir()+"/templates/find_news.gohtml",
	)
	if err != nil {
		panic(err)
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
		panic(err)
	}
}

func RenderUnsafe(s string) template.HTML {
	return template.HTML(s)
}
