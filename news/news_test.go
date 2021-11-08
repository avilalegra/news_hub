package news

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	kf := &KeeperFinderFake{}
	triggerA := make(chan time.Time)
	providerA := NewProviderMock(triggerA, Previews[0:2], nil)
	triggerB := make(chan time.Time)
	providerB := NewProviderMock(triggerB, Previews[2:], nil)

	collector := Collector{
		[]Provider{providerA, providerB},
		kf,
		log.Default(),
	}

	go collector.Run(context.Background())

	triggerA <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, Previews[:2], kf.Previews)

	triggerB <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, Previews, kf.Previews)
}

func TestCollectorContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	providerA := NewProviderMock(make(chan time.Time), Previews[0:2], nil)
	providerB := NewProviderMock(make(chan time.Time), Previews[2:], nil)
	exit := make(chan bool)
	collector := Collector{
		[]Provider{providerA, providerB},
		nil,
		log.Default(),
	}

	go func() {
		collector.Run(ctx)
		exit <- true
	}()
	<-time.After(1 * time.Millisecond)
	cancel()

	assert.True(t, <-exit)
	for _, provider := range collector.Providers {
		assert.Equal(t, ctx, provider.(*ProviderMock).Ctx)
	}
}

func TestCollectorLogsProvidersErrors(t *testing.T) {
	trigger := make(chan time.Time, 1)
	provider := NewProviderMock(trigger, nil, []error{errors.New("error fetching from source: rtve")})
	writerMock := new(writerMock)
	collector := Collector{
		[]Provider{provider},
		nil,
		log.New(writerMock, "", log.LstdFlags),
	}

	go collector.Run(context.Background())

	trigger <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Contains(t, writerMock.msg, "error fetching from source: rtve")
}

func TestCleaner(t *testing.T) {
	trigger := make(chan time.Time, 1)
	kf := KeeperFinderFake{Previews: tsCleaner}
	cleaner := Cleaner{&kf, trigger, int64((24 * time.Hour).Seconds())}

	go cleaner.Run(context.Background())
	trigger <- time.Now()
	close(trigger)
	<-time.After(1 * time.Millisecond)
	assert.Equal(t, 2, len(kf.Previews))
}

func TestCleanerContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	trigger := make(chan time.Time, 1)
	kf := KeeperFinderFake{Previews: tsCleaner}
	cleaner := Cleaner{&kf, trigger, int64((24 * time.Hour).Seconds())}

	exit := make(chan bool)

	go func() {
		cleaner.Run(ctx)
		exit <- true
	}()
	<-time.After(1 * time.Millisecond)
	cancel()

	assert.True(t, <-exit)
}

func TestPreviewMatchSearchExpr(t *testing.T) {
	for _, tData := range tsPreviewMatchSearchExpr {
		assert.Equal(t, tData.match, tData.preview.MatchPercent(tData.expr), tData.expr)
	}
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

var tsPreviewMatchSearchExpr = []struct {
	preview Preview
	expr    string
	match   int
}{
	{
		Preview{
			Title:       `Erupción en La Palma, en directo | La lava llega a 800 metros del mar y cambia de dirección al norte`,
			Description: `<ul> <li>Varios n&uacute;cleos poblacionales del municipio de Tazacorte han sido confinados</li> <li>La colada de lava podr&iacute;a llegar a la costa en las pr&oacute;ximas horas</li> </ul><br/><a href="http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml">Leer la noticia completa</a><img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&cg=F-N-B-TENOTICI-TESESPE01-TES800089&si=http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml" alt=""/>`,
		},
		"núcleos poblacionales la palma",
		100,
	},
	{
		Preview{
			Title:       "Energías renovables",
			Description: "La cumbre sobre el cambio climático, COP26, que se celebra en la ciudad de Glasgow (Escocia) ha incidido en que las energías fósiles son las más contaminantes para el medioambiente",
		},
		"cumbre Escocia.",
		100,
	},
	{
		Preview{
			Title:       "Energías renovables",
			Description: "La cumbre sobre el cambio climático, COP26, que se celebra en la ciudad de Glasgow (Escocia) ha incidido en que las energías fósiles son las más contaminantes para el medioambiente",
		},
		"Cumbre glasgow; energías renovables.",
		100,
	},
	{
		Preview{
			Title:       "Energías renovables",
			Description: "La cumbre sobre el cambio climático, COP26, que se celebra en la ciudad de Glasgow (Escocia) ha incidido en que las energías fósiles son las más contaminantes para el medioambiente",
		},
		"Noticias cumbre Glasgow energías renovables",
		80,
	},
	{
		Preview{
			Title:       "Energías renovables",
			Description: "La cumbre sobre el cambio climático, COP26, que se celebra en la ciudad de Glasgow (Escocia) ha incidido en que las energías fósiles son las más contaminantes para el medioambiente",
		},
		"energía renovable",
		100,
	},
	{
		Preview{
			Title:       "Energías renovables",
			Description: "La cumbre sobre el cambio climático, COP26, que se celebra en la ciudad de Glasgow (Escocia) ha incidido en que las energías fósiles son las más contaminantes para el medioambiente",
		},
		"combustibles fósiles",
		50,
	},
}
