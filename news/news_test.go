package news

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	Load(previews...)
	all := All()
	for _, tData := range previews {
		assert.Contains(t, all, tData)
	}
}

func TestSearch(t *testing.T) {
	Load(previews...)
	Load(previews...)
	for _, tData := range searchTestData {
		results := Search(tData.keywords)
		assert.Equal(t, tData.count, len(results), tData.keywords)
	}
}

type RepoMock struct {
	Previews []Preview
}

func (r *RepoMock) Add(preview Preview) {
	r.Previews = append(r.Previews, preview)
}

type ProviderMock struct {
	Trigger  chan time.Time
	Previews []Preview
	Errors   []error
}

func (p ProviderMock) RunAsync(providers chan<- Preview, errs chan<- error) {
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
	repo := &RepoMock{}
	triggerA := make(chan time.Time)
	providerA := ProviderMock{triggerA, previews[0:2], nil}
	triggerB := make(chan time.Time)
	providerB := ProviderMock{triggerB, previews[2:], nil}

	collector := Collector{
		[]Provider{providerA, providerB},
		repo,
		log.Default(),
	}
	collector.Run()

	triggerA <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, previews[:2], repo.Previews)

	triggerB <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, previews, repo.Previews)
}

func TestProviderErrorLog(t *testing.T) {
	repo := new(RepoMock)
	triggerA := make(chan time.Time, 1)
	providerA := ProviderMock{triggerA, nil, []error{errors.New("expected element type <rss> but have <xml>")}}
	triggerB := make(chan time.Time, 1)
	providerB := ProviderMock{triggerB, nil, []error{errors.New("bad server response when fetching xml")}}
	writerMock := new(WriterMock)
	logger := log.New(writerMock, "", log.LstdFlags)

	collector := Collector{
		[]Provider{providerA, providerB},
		repo,
		logger,
	}

	collector.Run()

	triggerA <- time.Now()
	time.Sleep(1 * time.Millisecond)
	assert.Contains(t, writerMock.msg, "expected element type <rss> but have <xml>")

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

var searchTestData = []struct {
	keywords string
	count    int
}{
	{
		"núcleos poblacionales",
		1,
	},
	{
		"Lava dirección confinados",
		1,
	},
	{
		"lava dirección confinados hierro",
		0,
	},
	{
		"directo, municipio",
		2,
	},
	{
		"Display; Support. PosTing",
		1,
	},
	{
		"<ul> <li>",
		0,
	},
	{
		"linux kernel",
		1,
	},
	{
		"linux kernel covid",
		0,
	},
}

var sources = map[string]*Source{
	"phoronix": {
		Title:       `Phoronix`,
		Link:        `https://www.phoronix.com/`,
		Language:    `en-US`,
		Description: `Linux Hardware Reviews & News`,
	},
	"rtve": {
		Title:       `Noticias en rtve.es`,
		Link:        `http://www.rtve.es`,
		Description: `RSS Tags`,
	},
}

var previews = []Preview{
	{
		Title:       `AMD Posts Code Enabling "Cyan Skillfish" Display Support Due To Different DCN2 Variant`,
		Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMD-Cyan-Skillfish-DCN-2.01`,
		Description: `Since July we've seen AMD open-source driver engineers posting code for "Cyan Skillfish" as an APU with Navi 1x graphics. While initial support for Cyan Skillfish was merged for Linux 5.15, it turns out the display code isn't yet wired up due to being a different DCN2 variant for its display block...`,
		Source:      sources["phoronix"],
	},
	{
		Title:       `Linux 5.16 To Bring Initial DisplayPort 2.0 Support For AMD Radeon Driver (AMDGPU)`,
		Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMDGPU-DP-2.0-Linux-5.16`,
		Description: `A batch of feature updates was submitted today for DRM-Next of early feature work slated to come to the next version of the Linux kernel...`,
		Source:      sources["phoronix"],
	},
	{
		Title:       `Erupción en La Palma, en directo | La lava llega a 800 metros del mar y cambia de dirección al norte`,
		Link:        `http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml`,
		Description: `<ul> <li>Varios n&uacute;cleos poblacionales del municipio de Tazacorte han sido confinados</li> <li>La colada de lava podr&iacute;a llegar a la costa en las pr&oacute;ximas horas</li> </ul><br/><a href="http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml">Leer la noticia completa</a><img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&cg=F-N-B-TENOTICI-TESESPE01-TES800089&si=http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml" alt=""/>`,
		Source:      sources["rtve"],
	},
	{
		Title:       `Guía de restricciones COVID: nuevas medidas en ocio nocturno, hostelería y aforos, directo`,
		Link:        `http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml`,
		Description: `<ul> <li>Repasa las principales medidas y restricciones frente a la COVID-19, comunidad a comunidad del municipio</li> <li><a href="https://www.rtve.es/noticias/20210928/coronavirus-covid-directo-espana-mundo-ultima-hora/2175601.shtml" target="_blank">Coronavirus: &uacute;ltima hora</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/mapa-del-coronavirus-espana/2004681.shtml" target="_blank">Mapa de Espa&ntilde;a</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/ocupacion-camas-covid-19-hospitales-espanoles/2042349.shtml" target="_blank">Hospitales y UCI</a></li> <li><a href="https://www.rtve.es/noticias/20210924/campana-vacunacion-espana/2062499.shtml" target="_blank">Vacunas en Espa&ntilde;a</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/mapa-mundial-del-coronavirus/1998143.shtml" target="_blank">Mapa mundial&#8203;</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/lab/vacunacion-espana-coronavirus/">Especial: La gran vacunaci&oacute;n</a></li> </ul><br/><a href="http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml">Leer la noticia completa</a><img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&cg=F-N-B-TENOTICI-TESESPE01-TELCO20VX&si=http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml" alt=""/>`,
		Source:      sources["rtve"],
	},
}
