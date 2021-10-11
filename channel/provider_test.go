package channel

import (
	"avilego.me/news_hub/news"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRssProvider(t *testing.T) {
	for _, tData := range tsRssProvider {
		trigger := make(chan time.Time, 2)
		previewsChan := make(chan news.Preview, 10)
		errorsChan := make(chan error, 10)
		var previews []news.Preview
		var errs []error
		provider := NewRssNewsProvider(tData.sources, trigger)

		provider.RunAsync(previewsChan, errorsChan)

		trigger <- time.Now()
		trigger <- time.Now()
		close(trigger)

		for preview := range previewsChan {
			previews = append(previews, preview)
		}
		for err := range errorsChan {
			errs = append(errs, err)
		}

		assert.Equal(t, len(tData.previews)*2, len(previews))
		for _, preview := range previews {
			assert.Contains(t, tData.previews, preview)
		}

		assert.Equal(t, len(tData.errors)*2, len(errs))
	}
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
