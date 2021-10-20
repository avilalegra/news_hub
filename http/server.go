package http

import (
	"avilego.me/recent_news/news"
)

func newPreviewData(preview news.Preview) previewData {
	return previewData{
		Title:       preview.Title,
		Link:        preview.Link,
		Description: preview.Description,
		SourceLink:  preview.Source.Link,
	}
}

type previewData struct {
	Title       string
	Link        string
	Description string
	SourceLink  string
}
