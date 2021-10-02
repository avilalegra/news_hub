package http

import (
	"encoding/json"
	"net/http"

	"avilego.me/news_hub/news"
)

func Search(w http.ResponseWriter, req *http.Request) {

	news.Load(news.PreviewsFakes())

	news := news.Search("something")

	jres, err := json.Marshal(news)

	if err != nil {
		panic("bad json")
	}

	w.Write([]byte(jres))
	w.WriteHeader(http.StatusBadRequest)
}

/*
func Run(addr string) {
	http.HandleFunc("/search", search)

	http.ListenAndServe(addr, nil)
}
*/
