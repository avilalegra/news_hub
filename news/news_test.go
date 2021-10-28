package news

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	writerMock := new(writerMock)
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
	writerMock := new(writerMock)
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

type writerMock struct {
	msg string
}

func (w *writerMock) Write(p []byte) (n int, err error) {
	w.msg = string(p)
	return 1, nil
}
