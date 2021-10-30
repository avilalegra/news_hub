package news

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	r := &KeeperFinderFake{}
	triggerA := make(chan time.Time)
	providerA := ProviderMock{triggerA, Previews[0:2], nil}
	triggerB := make(chan time.Time)
	providerB := ProviderMock{triggerB, Previews[2:], nil}

	collector := Collector{
		[]Provider{providerA, providerB},
		r,
		log.Default(),
	}

	go collector.Run()

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
		[]Provider{provider},
		nil,
		log.New(writerMock, "", log.LstdFlags),
	}

	go collector.Run()

	trigger <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Contains(t, writerMock.msg, "error fetching from source: rtve")
}

func TestCleaner(t *testing.T) {
	trigger := make(chan time.Time, 1)
	kf := KeeperFinderFake{Previews: tsCleaner}
	cleaner := Cleaner{&kf, trigger, int64((24 * time.Hour).Seconds())}

	cleaner.Run()
	trigger <- time.Now()
	close(trigger)
	<-time.After(1 * time.Millisecond)
	assert.Equal(t, 2, len(kf.Previews))
}

type writerMock struct {
	msg string
}

func (w *writerMock) Write(p []byte) (n int, err error) {
	w.msg = string(p)
	return 1, nil
}

var tsCleaner = []Preview{
	{
		Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMDGPU-DP-2.0-Linux-5.16`,
		RegUnixTime: time.Now().Unix() - int64((25 * time.Hour).Seconds()),
	},
	{
		Link:        `http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml`,
		RegUnixTime: time.Now().Unix(),
	},
	{
		Link:        `http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml`,
		RegUnixTime: time.Now().Unix() - int64((24 * time.Hour).Seconds()) + 10,
	},
}
