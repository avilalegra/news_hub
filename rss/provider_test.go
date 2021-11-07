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

func TestRssProviderSendPreviews(t *testing.T) {
	for i, tData := range tsProviderSendPreviews {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			t.Parallel()

			previews, errs := collectProvider(tData.sources)

			assert.Empty(t, errs)
			for _, p := range tData.previews {
				assert.Contains(t, previews, p)
			}
		})
	}
}

func TestRssProviderSendErrors(t *testing.T) {
	for i, tData := range tsProviderSendErrors {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			t.Parallel()

			previews, errs := collectProvider(tData.sources)

			assert.Empty(t, previews)
			assert.Equal(t, len(tData.errors), len(errs))
		})
	}
}

func TestProviderContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	provider := NewRssProvider(nil, make(chan time.Time))
	exit := make(chan bool)

	go func() {
		provider.Provide(ctx, nil, nil)
		exit <- true
	}()

	cancel()
	ok := <-exit
	assert.True(t, ok)
}

func collectProvider(sources []Source) (previews []news.Preview, errors []error) {
	previewsChan := make(chan news.Preview)
	errorsChan := make(chan error)
	waitProvider := make(chan bool)
	trigger := make(chan time.Time)
	ctx, cancel := context.WithCancel(context.TODO())

	provider := NewRssProvider(sources, trigger)

	go func() {
		for p := range previewsChan {
			previews = append(previews, p)
		}
	}()
	go func() {
		for e := range errorsChan {
			errors = append(errors, e)
		}
	}()

	go func() {
		provider.Provide(ctx, previewsChan, errorsChan)
		waitProvider <- true
	}()

	trigger <- time.Now()
	cancel()
	<-waitProvider

	return previews, errors
}

var tsProviderSendPreviews = []struct {
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
		[]Source{},
		[]news.Preview{},
	},
}

var tsProviderSendErrors = []struct {
	sources []Source
	errors  []error
}{
	{
		[]Source{
			{"http://sample/url/2", newHttpClientMock("http://sample/url/2", `<?xml version="1.0"?><xml></xl>`)},
		},
		[]error{errors.New("expected element type <rss> but have <xml>")},
	},
}
