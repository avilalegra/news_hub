package http

import (
	"avilego.me/recent_news/news"
)

type searchResponse struct {
	Count int
	Data  searchData
}

type searchData struct {
	Sources  map[string]news.Source
	Previews []previewData
}

type previewData struct {
	Title       string
	Link        string
	Description string
	SourceLink  string
}

func newPreviewData(preview news.Preview) previewData {
	return previewData{
		Title:       preview.Title,
		Link:        preview.Link,
		Description: preview.Description,
		SourceLink:  preview.Source.Link,
	}
}

func newSearchResponse(previews []news.Preview) searchResponse {
	sources := make(map[string]news.Source)
	prvsData := make([]previewData, len(previews))

	for i, preview := range previews {
		sources[preview.Source.Link] = *preview.Source
		prvsData[i] = newPreviewData(preview)
	}

	return searchResponse{
		Count: len(prvsData),
		Data: searchData{
			Sources:  sources,
			Previews: prvsData,
		},
	}
}
