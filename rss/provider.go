package rss

import (
	"context"
	"sync"
	"time"

	"avilego.me/recent_news/news"
)

func NewRssProvider(sources []Source, interval <-chan time.Time) news.Provider {
	return Provider{
		sources,
		interval,
	}
}

type Provider struct {
	sources  []Source
	interval <-chan time.Time
}

func (p Provider) Provide(ctx context.Context, previewsChan chan<- news.Preview, errorsChan chan<- error) {
	for running := true; running; {
		select {
		case <-p.interval:
			p.fetchSourcesNews(previewsChan, errorsChan)
		case <-ctx.Done():
			running = false
		}
	}
}

func (p Provider) fetchSourcesNews(previewsChan chan<- news.Preview, errorsChan chan<- error) {
	var wg sync.WaitGroup
	for _, source := range p.sources {
		wg.Add(1)
		go func(s Source) {
			defer wg.Done()

			if channel, err := s.Fetch(); err == nil {
				for _, preview := range channel.GetNews() {
					previewsChan <- preview
				}
			} else {
				errorsChan <- err
			}

		}(source)
	}
	wg.Wait()
}
