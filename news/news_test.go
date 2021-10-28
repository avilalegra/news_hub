package news

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type KeeperMock struct {
	Previews []Preview
	Error    error
}

func (r *KeeperMock) Store(preview Preview) error {
	if r.Error != nil {
		return r.Error
	}
	r.Previews = append(r.Previews, preview)
	return nil
}

func (r *KeeperMock) Remove(preview Preview) {
	panic("implement me")
}

func NewMockKeeper() *KeeperMock {
	return &KeeperMock{make([]Preview, 0), nil}
}

func NewFailingMockKeeper(err error) *KeeperMock {
	return &KeeperMock{make([]Preview, 0), err}
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
	providerA := ProviderMock{triggerA, Previews[0:2], nil}
	triggerB := make(chan time.Time)
	providerB := ProviderMock{triggerB, Previews[2:], nil}

	collector := Collector{
		[]AsyncProvider{providerA, providerB},
		r,
		log.Default(),
	}
	collector.Run()

	triggerA <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, Previews[:2], r.Previews)

	triggerB <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, Previews, r.Previews)
}

func TestCollectorLogsProvidersErrors(t *testing.T) {
	trigger := make(chan time.Time, 1)
	provider := ProviderMock{trigger, nil, []error{errors.New("error fetching from source: rtve")}}
	writerMock := new(WriterMock)
	collector := Collector{
		[]AsyncProvider{provider},
		NewMockKeeper(),
		log.New(writerMock, "", log.LstdFlags),
	}

	collector.Run()

	trigger <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Contains(t, writerMock.msg, "error fetching from source: rtve")
}

func TestCollectorLogsKeeperErrors(t *testing.T) {
	keeperErr := errors.New("couldn't save preview")
	trigger := make(chan time.Time, 1)
	provider := ProviderMock{trigger, Previews, nil}
	writerMock := new(WriterMock)
	collector := Collector{
		[]AsyncProvider{provider},
		NewFailingMockKeeper(keeperErr),
		log.New(writerMock, "", log.LstdFlags),
	}

	collector.Run()

	trigger <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Contains(t, writerMock.msg, "couldn't save preview")
}

type WriterMock struct {
	msg string
}

func (w *WriterMock) Write(p []byte) (n int, err error) {
	w.msg = string(p)
	return 1, nil
}
