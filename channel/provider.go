package channel

import (
	"sync"
	"time"

	"avilego.me/news_hub/news"
)

func NewRssNewsProvider(sources []Source, interval chan time.Time) news.Provider {
	return RssNewsProvider{
		sources,
		interval,
	}
}

type RssNewsProvider struct {
	sources  []Source
	interval chan time.Time
}

func (p RssNewsProvider) RunAsync(previewsChan chan<- news.Preview, errorsChan chan<- error) {
	go func() {
		for range p.interval {
			var wg sync.WaitGroup
			for _, source := range p.sources {
				wg.Add(1)
				go func(s Source) {
					defer wg.Done()
					fetchSourceNews(s, previewsChan, errorsChan)
				}(source)
			}
			wg.Wait()
		}
		close(previewsChan)
		close(errorsChan)
	}()
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
