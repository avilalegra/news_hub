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

func (p RssNewsProvider) RunAsync(previewsChan chan<- news.Preview) {
	go func() {
		var wg sync.WaitGroup
		for range p.interval {
			for _, source := range p.sources {
				wg.Add(1)
				go func(s Source) {
					defer wg.Done()
					previews, _ := s.FetchNews()
					for _, preview := range previews {
						previewsChan <- preview
					}
				}(source)
			}
		}
		wg.Wait()
		close(previewsChan)
	}()
}
