package rss

import (
	"avilego.me/recent_news/news"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//TODO: Simplify this test
func TestRssProvider(t *testing.T) {
	for i, tData := range tsRssProvider {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			t.Parallel()
			trigger := make(chan time.Time, 2)
			previewsChan := make(chan news.Preview, 10)
			errorsChan := make(chan error, 10)
			var previews []news.Preview
			var errs []error
			provider := NewRssNewsProvider(tData.sources, trigger)

			go provider.Provide(previewsChan, errorsChan, context.TODO())

			trigger <- time.Now()
			trigger <- time.Now()
			close(trigger)

			for i := 0; i < len(tData.previews)*2; i++ {
				previews = append(previews, <-previewsChan)
			}

			for i := 0; i < len(tData.errors)*2; i++ {
				errs = append(errs, <-errorsChan)
			}

			assert.Equal(t, len(tData.previews)*2, len(previews))
			for _, preview := range previews {
				assert.Contains(t, tData.previews, preview)
			}

			assert.Equal(t, len(tData.errors)*2, len(errs))
		})
	}
}

func TestProviderContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	provider := NewRssNewsProvider(nil, make(chan time.Time))
	exit := make(chan bool)

	go func() {
		provider.Provide(nil, nil, ctx)
		exit <- true
	}()

	cancel()
	ok := <-exit
	assert.True(t, ok)
}

var tsRssProvider = []struct {
	sources  []Source
	previews []news.Preview
	errors   []error
}{
	{
		[]Source{
			{"http://sample/url/1", newHttpClientMock("http://sample/url/1", chanSamples[3].xml)},
			{"http://sample/url/2", newHttpClientMock("http://sample/url/2", chanSamples[4].xml)},
		},
		append(append(make([]news.Preview, 0), chanSamples[3].channel.GetNews()...), chanSamples[4].channel.GetNews()...),
		nil,
	},

	{
		[]Source{
			{"http://sample/url/2", newHttpClientMock("http://sample/url/2", chanSamples[4].xml)},
		},
		append(make([]news.Preview, 0), chanSamples[4].channel.GetNews()...),
		nil,
	},
	{
		[]Source{
			{"http://sample/url/1", newHttpClientMock("http://sample/url/1", chanSamples[3].xml)},
			{"http://sample/url/2", newHttpClientMock("http://sample/url/2", `<?xml version="1.0"?><xml></xl>`)},
		},
		chanSamples[3].channel.GetNews(),
		[]error{errors.New("expected element type <rss> but have <xml>")},
	},
	{
		[]Source{},
		[]news.Preview{},
		nil,
	},
}
