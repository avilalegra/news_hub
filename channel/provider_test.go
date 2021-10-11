package channel

import (
	"avilego.me/news_hub/news"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRssProvider(t *testing.T) {
	for _, tData := range tsRssProvider {
		trigger := make(chan time.Time, 2)
		previewsChan := make(chan news.Preview, 10)
		var previews []news.Preview
		provider := NewRssNewsProvider(tData.sources, trigger)

		provider.RunAsync(previewsChan)

		trigger <- time.Now()
		trigger <- time.Now()
		close(trigger)

		for preview := range previewsChan {
			previews = append(previews, preview)
		}

		assert.Equal(t, len(tData.previews)*2, len(previews))
		for _, preview := range previews {
			assert.Contains(t, tData.previews, preview)
		}
	}
}

var tsRssProvider = []struct {
	sources  []Source
	previews []news.Preview
}{
	{
		[]Source{
			{"http://sample/url/1", newHttpClientMock("http://sample/url/1", chanSamples[3].xml)},
			{"http://sample/url/2", newHttpClientMock("http://sample/url/2", chanSamples[4].xml)},
		},
		append(append(make([]news.Preview, 0), chanSamples[3].channel.GetNews()...), chanSamples[4].channel.GetNews()...),
	},

	{
		[]Source{
			{"http://sample/url/2", newHttpClientMock("http://sample/url/2", chanSamples[4].xml)},
		},
		append(make([]news.Preview, 0), chanSamples[4].channel.GetNews()...),
	},
	{
		[]Source{
			{"http://sample/url/1", newHttpClientMock("http://sample/url/1", chanSamples[3].xml)},
			{"http://sample/url/2", newHttpClientMock("http://sample/url/2", `<?xml version="1.0"?><xml></xl>`)},
		},
		chanSamples[3].channel.GetNews(),
	},
	{
		[]Source{},
		[]news.Preview{},
	},
}
