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
				go func(src Source) {
					channel, err := src.Fetch()
					if err == nil {
						for _, preview := range channel.GetNews() {
							previewsChan <- preview
						}
					}
					wg.Done()
				}(source)
			}
		}
		wg.Wait()
		close(previewsChan)
	}()
}
