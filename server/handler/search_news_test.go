package handler

import (
	"avilego.me/recent_news/news"
	"avilego.me/recent_news/news/newstest"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestNewPreviewData(t *testing.T) {
	preview := newstest.Previews[0]
	prvData := newPreviewData(preview)
	assert.Equal(t, previewData{
		Title:       preview.Title,
		Link:        preview.Link,
		Description: preview.Description,
		SourceLink:  preview.Source.Link,
	}, prvData)
}

func TestNewSearchResponse(t *testing.T) {
	for _, tData := range tsMakeSearchResponse {
		response := newSearchResponse(tData.previews)
		assert.Equal(t, tData.response, response)
	}
}

func TestSearch(t *testing.T) {
	for _, tData := range tsSearch {
		handler := SearchHandler{newstest.Finder{Keywords: tData.keywords, Previews: tData.previews}}
		params := url.Values{}
		params.Set("keywords", tData.keywords)

		expectedJson, _ := json.Marshal(newSearchResponse(tData.previews))

		assert.HTTPBodyContains(t, handler.ServeHTTP, "GET", "/search", params, string(expectedJson))
	}
}

var tsMakeSearchResponse = []struct {
	previews []news.Preview
	response searchResponse
}{
	{
		newstest.Previews[:2],
		searchResponse{
			Count: 2,
			Data: searchData{
				Sources: []news.Source{
					*newstest.Sources["phoronix"],
				},
				Previews: []previewData{
					newPreviewData(newstest.Previews[0]),
					newPreviewData(newstest.Previews[1]),
				},
			},
		},
	},
	{
		newstest.Previews[:3],
		searchResponse{
			Count: 3,
			Data: searchData{
				Sources: []news.Source{
					*newstest.Sources["phoronix"],
					*newstest.Sources["rtve"],
				},
				Previews: []previewData{
					newPreviewData(newstest.Previews[0]),
					newPreviewData(newstest.Previews[1]),
					newPreviewData(newstest.Previews[2]),
				},
			},
		},
	},
	{
		newstest.Previews[0:0],
		searchResponse{
			Count: 0,
			Data: searchData{
				Sources:  []news.Source{},
				Previews: []previewData{},
			},
		},
	},
}

var tsSearch = []struct {
	keywords string
	previews []news.Preview
}{
	{
		"AMD",
		newstest.Previews[0:2],
	},
	{
		"amd",
		newstest.Previews[0:2],
	},
	{
		"",
		newstest.Previews[0:0],
	},
}
