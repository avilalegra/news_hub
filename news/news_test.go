package news

import (
	"avilego.me/recent_news/news/newstest"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type KeeperMock struct {
	Previews []Preview
}

func (r *KeeperMock) Store(preview Preview) error {
	r.Previews = append(r.Previews, preview)
	return nil
}

type ProviderMock struct {
	Trigger  chan time.Time
	Previews []Preview
	Errors   []error
}

func (p ProviderMock) ProvideAsync(providers chan<- Preview, errs chan<- error) {
	go func() {
		for range p.Trigger {
			for _, preview := range p.Previews {
				providers <- preview
			}
			for _, e := range p.Errors {
				errs <- e
			}
		}
	}()
}

func TestCollector(t *testing.T) {
	r := &KeeperMock{}
	triggerA := make(chan time.Time)
	providerA := ProviderMock{triggerA, newstest.Previews[0:2], nil}
	triggerB := make(chan time.Time)
	providerB := ProviderMock{triggerB, newstest.Previews[2:], nil}

	collector := Collector{
		[]AsyncProvider{providerA, providerB},
		r,
		log.Default(),
	}
	collector.Run()

	triggerA <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, newstest.Previews[:2], r.Previews)

	triggerB <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, newstest.Previews, r.Previews)
}

func TestProviderErrorLog(t *testing.T) {
	r := new(KeeperMock)
	triggerA := make(chan time.Time, 1)
	providerA := ProviderMock{triggerA, newstest.Previews[0:1], nil}
	triggerB := make(chan time.Time, 1)
	providerB := ProviderMock{triggerB, nil, []error{errors.New("bad server response when fetching xml")}}
	writerMock := new(WriterMock)
	logger := log.New(writerMock, "", log.LstdFlags)

	collector := Collector{
		[]AsyncProvider{providerA, providerB},
		r,
		logger,
	}

	collector.Run()

	triggerA <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Contains(t, writerMock.msg, `news preview added: AMD Posts Code Enabling "Cyan Skillfish" Display Support Due To Different DCN2 Variant`)

	triggerB <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Contains(t, writerMock.msg, "bad server response when fetching xml")
}

type WriterMock struct {
	msg string
}

func (w *WriterMock) Write(p []byte) (n int, err error) {
	w.msg = string(p)
	return 1, nil
}
