package api

import (
	"avilego.me/recent_news/news"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPreviewData(t *testing.T) {
	preview := news.Previews[0]
	prvData := newPreviewData(preview)
	assert.Equal(t, previewData{
		Title:       preview.Title,
		Link:        preview.Link,
		Description: preview.Description,
		SourceLink:  preview.Source.Link,
	}, prvData)
}

func TestNewSearchResponse(t *testing.T) {
	for i, tData := range tsMakeSearchResponse {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			t.Parallel()
			response := newSearchResponse(tData.previews)
			assert.Equal(t, tData.response, response)
		})
	}
}

func TestSearch(t *testing.T) {
	for i, tData := range tsSearch {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			t.Parallel()
			handler := SearchHandler{news.KeeperFinderFake{Previews: tData.previews}}
			expectedJson, _ := json.Marshal(newSearchResponse(tData.previews))
			request, _ := http.NewRequest("GET", "/api/search?keywords="+tData.keywords, nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, request)

			assert.Equal(t, 200, resp.Code)
			assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
			assert.Contains(t, string(expectedJson), resp.Body.String())
		})
	}
}

var tsMakeSearchResponse = []struct {
	previews []news.Preview
	response searchResponse
}{
	{
		news.Previews[:2],
		searchResponse{
			Count: 2,
			Data: searchData{
				Sources: []news.Source{
					*news.Sources["phoronix"],
				},
				Previews: []previewData{
					newPreviewData(news.Previews[0]),
					newPreviewData(news.Previews[1]),
				},
			},
		},
	},
	{
		news.Previews[:3],
		searchResponse{
			Count: 3,
			Data: searchData{
				Sources: []news.Source{
					*news.Sources["phoronix"],
					*news.Sources["rtve"],
				},
				Previews: []previewData{
					newPreviewData(news.Previews[0]),
					newPreviewData(news.Previews[1]),
					newPreviewData(news.Previews[2]),
				},
			},
		},
	},
	{
		news.Previews[0:0],
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
		news.Previews[0:2],
	},
	{
		"amd",
		news.Previews[0:2],
	},
	{
		"",
		news.Previews[0:0],
	},
}
