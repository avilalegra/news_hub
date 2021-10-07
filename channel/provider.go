package channel

import (
	"log"
	"sync"
	"time"

	"avilego.me/news_hub/news"
)

func NewRssNewsProvider(sources []Source, interval chan time.Time, logger *log.Logger) news.Provider {
	return RssNewsProvider{
		sources,
		interval,
		logger,
	}
}

type RssNewsProvider struct {
	sources  []Source
	interval chan time.Time
	logger   *log.Logger
}

func (p RssNewsProvider) RunAsync(previewsChan chan<- news.Preview) {
	go func() {
		var wg sync.WaitGroup
		for range p.interval {
			for _, source := range p.sources {
				wg.Add(1)
				go func(s Source) {
					defer wg.Done()
					previews, err := s.FetchNews()
					if err != nil {
						p.logger.Println(err)
					}
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
