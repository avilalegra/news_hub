package rss

import (
	"context"
	"sync"
	"time"

	"avilego.me/recent_news/news"
)

func NewRssNewsProvider(sources []Source, interval <-chan time.Time) news.Provider {
	return NewsProvider{
		sources,
		interval,
	}
}

type NewsProvider struct {
	sources  []Source
	interval <-chan time.Time
}

func (p NewsProvider) Provide(ctx context.Context, previewsChan chan<- news.Preview, errorsChan chan<- error) {
	for running := true; running; {
		select {
		case <-p.interval:
			var wg sync.WaitGroup
			for _, source := range p.sources {
				wg.Add(1)
				go func(s Source) {
					defer wg.Done()
					fetchSourceNews(s, previewsChan, errorsChan)
				}(source)
			}
			wg.Wait()
		case <-ctx.Done():
			running = false
		}
	}
}

func fetchSourceNews(s Source, previewsChan chan<- news.Preview, errorsChan chan<- error) {
	if channel, err := s.Fetch(); err == nil {
		for _, preview := range channel.GetNews() {
			previewsChan <- preview
		}
	} else {
		errorsChan <- err
	}
}
