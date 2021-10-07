package channel

import (
	"avilego.me/news_hub/news"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestRssProvider(t *testing.T) {
	for _, tData := range tsRssProvider {
		trigger := make(chan time.Time, 2)
		previewsChan := make(chan news.Preview, 10)
		var previews []news.Preview
		provider := NewRssNewsProvider(tData.sources, trigger, log.Default())

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

func TestRssProviderErrorLog(t *testing.T) {
	sources := []Source{{"http://sample/url/2", newHttpClientMock("http://sample/url/2", `<?xml version="1.0"?><xml></xl>`)}}
	previewsChan := make(chan news.Preview, 1)
	trigger := make(chan time.Time, 1)
	writerMock := new(WriterMock)
	logger := log.New(writerMock, "", log.LstdFlags)
	provider := NewRssNewsProvider(sources, trigger, logger)

	provider.RunAsync(previewsChan)
	trigger <- time.Now()
	close(trigger)
	<-previewsChan

	assert.Contains(t, writerMock.msg, "expected element type <rss> but have <xml>")
}

type WriterMock struct {
	msg string
}

func (w *WriterMock) Write(p []byte) (n int, err error) {
	w.msg = string(p)
	return 1, nil
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
