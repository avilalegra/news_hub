package channel

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRssParsing(t *testing.T) {
	for _, tData := range tsParseChannelTests() {
		result, e := Parse([]byte(tData.xml))

		switch v := tData.result.(type) {
		case *Channel:
			assertChannelEquals(t, *v, *result)
		case error:
			assert.Error(t, e)
		}
	}
}

func TestFetchChannel(t *testing.T) {
	for _, tData := range tsFetchChannel() {
		result, e := tData.source.Fetch()

		switch v := tData.result.(type) {
		case *Channel:
			assertChannelEquals(t, *v, *result)
		case error:
			assert.Error(t, e)
		}
	}
}

func TestGetChannelNews(t *testing.T) {
	for _, channel := range channelNewsTests {
		for i, chanNews := range channel.GetNews() {
			assert.Equal(t, channel.Title, chanNews.Source.Title)
			assert.Equal(t, channel.Link, chanNews.Source.Link)
			assert.Equal(t, channel.Language, chanNews.Source.Language)
			assert.Equal(t, channel.Description, chanNews.Source.Description)

			assert.Equal(t, channel.Items[i].Title, chanNews.Title)
			assert.Equal(t, channel.Items[i].Link, chanNews.Link)
			assert.Equal(t, channel.Items[i].Description, chanNews.Description)
		}
	}
}

func newHttpClientMock(link string, response interface{}) HttpClient {
	client := new(HttpClientMock)
	switch v := response.(type) {
	case error:
		client.On("Get", link).Return([]byte(""), v)
	case string:
		client.On("Get", link).Return([]byte(v), nil)
	}
	return client
}

func assertChannelEquals(t *testing.T, expected Channel, actual Channel) {
	assert.Equal(t, expected.Title, actual.Title, "Title parsing error")
	assert.Equal(t, expected.Link, actual.Link, "Link parsing error")
	assert.Equal(t, expected.Description, actual.Description, "Description parsing error")
	assert.Equal(t, expected.Language, actual.Language, "Language parsing error")
	assert.Equal(t, expected.LastBuildDate, actual.LastBuildDate, "LastBuildDate parsing error")

	for i, expitem := range expected.Items {
		assert.Equal(t, expitem.Title, actual.Items[i].Title, "Item title parsing error")
		assert.Equal(t, expitem.Link, actual.Items[i].Link, "Item link parsing error")
		assert.Equal(t, expitem.Description, actual.Items[i].Description, "Item description parsing error")
	}
}

type HttpClientMock struct {
	mock.Mock
}

func (m *HttpClientMock) Get(url string) ([]byte, error) {
	args := m.Called(url)
	return args.Get(0).([]byte), args.Error(1)
}

func parseTime(timeExpr string) *time.Time {
	ptime, _ := time.Parse(time.RFC1123, timeExpr)
	return &ptime
}

type chanSample struct {
	xml     string
	channel Channel
}

type fetchChannelSample struct {
	source Source
	result interface{}
}

func tsFetchChannel() []fetchChannelSample {
	samples := make([]fetchChannelSample, len(chanSamples), len(chanSamples)+2)

	for i := 0; i < len(chanSamples); i++ {
		url := fmt.Sprintf("http://sample/url/%d", i)
		samples[i] = fetchChannelSample{
			source: Source{url, newHttpClientMock(url, chanSamples[i].xml)},
			result: chanSamples[i].channel,
		}
	}

	return append(samples,
		fetchChannelSample{
			source: Source{"http://sample/url/error/1", newHttpClientMock("http://sample/url/error/1", errors.New("Parsing error"))},
			result: errors.New("Parsing error"),
		},
		fetchChannelSample{
			source: Source{"http://sample/url/error/2", newHttpClientMock("http://sample/url/error/2", errors.New("500 http response"))},
			result: errors.New("500 http response"),
		})
}

type parseChanSample struct {
	xml    string
	result interface{}
}

func tsParseChannelTests() []parseChanSample {
	samples := make([]parseChanSample, len(chanSamples), len(chanSamples)+3)
	for i := 0; i < len(chanSamples); i++ {
		samples[i] = parseChanSample{chanSamples[i].xml, chanSamples[i].channel}
	}
	return append(samples,
		parseChanSample{
			`<?xml version="1.0"?><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/"><channel rdf:about="http://www.xml.com/xml/news.rss"><title>XML.com</title><link>http://xml.com/pub</link><description>rss v1 description
	</description><image rdf:resource="http://xml.com/universal/images/xml_tiny.gif" /><items><rdf:Seq><rdf:li resource="http://xml.com/pub/2000/08/09/xslt/xslt.html" /><rdf:li resource="http://xml.com/pub/2000/08/09/rdfdb/index.html" /></rdf:Seq></items><textinput rdf:resource="http://search.xml.com" /></channel><image rdf:about="http://xml.com/universal/images/xml_tiny.gif"><title>XML.com</title><link>http://www.xml.com</link><url>http://xml.com/universal/images/xml_tiny.gif</url></image><item rdf:about="http://xml.com/pub/2000/08/09/xslt/xslt.html"><title>Processing Inclusions with XSLT</title><link>http://xml.com/pub/2000/08/09/xslt/xslt.html</link><description>rss description
	</description></item><item rdf:about="http://xml.com/pub/2000/08/09/rdfdb/index.html"><title>Putting RDF to Work</title><link>http://xml.com/pub/2000/08/09/rdfdb/index.html</link><description>item description</description></item><textinput rdf:about="http://search.xml.com"><title>Search XML.com</title><description>Search XML.com's XML collection</description><name>s</name><link>http://search.xml.com</link></textinput></rdf:RDF>`,
			errors.New("Invalid xml 1"),
		},
		parseChanSample{
			`<?xml version="1.0"?><xml></xl>`,
			errors.New("Invalid xml 1"),
		},
		parseChanSample{
			"",
			errors.New("Invalid xml 2"),
		})
}

var channelNewsTests = []Channel{
	chanSamples[3].channel,
	chanSamples[4].channel,
}

var chanSamples = []chanSample{
	{
		`<rss version="2.0"><channel><title>Phoronix</title><link>https://www.phoronix.com/</link><description><![CDATA[Linux Hardware Reviews & News]]></description></channel></rss>`,
		Channel{
			Title:       `Phoronix`,
			Link:        `https://www.phoronix.com/`,
			Description: `Linux Hardware Reviews & News`,
		},
	},
	{
		`<?xml version="1.0" encoding="UTF-8"?><rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:media="http://search.yahoo.com/mrss/" xmlns:nyt="http://www.nytimes.com/namespaces/rss/2.0" version="2.0"><channel><title>NYT &gt; Top Stories</title><link>https://www.nytimes.com</link><description>NYT channel description</description><language>en-us</language><lastBuildDate>Sun, 19 Sep 2021 06:27:36 +0000</lastBuildDate><image><title>NYT > Top Stories</title><url>https://static01.nyt.com/images/misc/NYT_logo_rss_250x40.png</url><link>https://www.nytimes.com</link></image></channel></rss>`,
		Channel{
			Title:         `NYT > Top Stories`,
			Link:          `https://www.nytimes.com`,
			Description:   `NYT channel description`,
			Language:      "en-us",
			LastBuildDate: Time{Time: parseTime(`Sun, 19 Sep 2021 06:27:36 +0000`)},
		},
	},
	{
		`<?xml version="1.0" encoding="UTF-8"?><rss xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xmlns:feedburner="http://rssnamespace.org/feedburner/ext/1.0" version="2.0"><channel><title><![CDATA[CNN.com - RSS Channel - World]]></title><description><![CDATA[CNN.com delivers up-to-the-minute news and information on the latest top stories, weather, entertainment, politics and more.]]></description><link>https://www.cnn.com/world/index.html</link><image><url>http://i2.cdn.turner.com/cnn/2015/images/09/24/cnn.digital.png</url><title>CNN.com - RSS Channel - World</title><link>https://www.cnn.com/world/index.html</link></image><lastBuildDate>Mon, 27 Sep 2021 13:54:35 GMT</lastBuildDate><pubDate>Thu, 16 Sep 2021 15:14:25 GMT</pubDate><language><![CDATA[en-US]]></language><atom10:link xmlns:atom10="http://www.w3.org/2005/Atom" rel="self" type="application/rss+xml" href="http://rss.cnn.com/rss/edition_world" /></channel></rss>`,
		Channel{
			Title:         `CNN.com - RSS Channel - World`,
			Link:          `https://www.cnn.com/world/index.html`,
			Language:      `en-US`,
			Description:   `CNN.com delivers up-to-the-minute news and information on the latest top stories, weather, entertainment, politics and more.`,
			LastBuildDate: Time{Time: parseTime(`Mon, 27 Sep 2021 13:54:35 GMT`)},
		},
	},
	{
		`<rss version="2.0"><channel><title>Phoronix</title><link>https://www.phoronix.com/</link><description>Linux Hardware Reviews & News</description><language>en-us</language><item><title>AMD Posts Code Enabling "Cyan Skillfish" Display Support Due To Different DCN2 Variant</title><link>https://www.phoronix.com/scan.php?page=news_item&px=AMD-Cyan-Skillfish-DCN-2.01</link><guid>https://www.phoronix.com/scan.php?page=news_item&px=AMD-Cyan-Skillfish-DCN-2.01</guid><description>Since July we've seen AMD open-source driver engineers posting code for "Cyan Skillfish" as an APU with Navi 1x graphics. While initial support for Cyan Skillfish was merged for Linux 5.15, it turns out the display code isn't yet wired up due to being a different DCN2 variant for its display block...</description><pubDate>Tue, 28 Sep 2021 00:00:00 -0400</pubDate></item><item><title>Linux 5.16 To Bring Initial DisplayPort 2.0 Support For AMD Radeon Driver (AMDGPU)</title><link>https://www.phoronix.com/scan.php?page=news_item&px=AMDGPU-DP-2.0-Linux-5.16</link><guid>https://www.phoronix.com/scan.php?page=news_item&px=AMDGPU-DP-2.0-Linux-5.16</guid><description>A batch of feature updates was submitted today for DRM-Next of early feature work slated to come to the next version of the Linux kernel...</description><pubDate>Mon, 27 Sep 2021 17:46:34 -0400</pubDate></item></channel></rss>`,
		Channel{
			Title:       `Phoronix`,
			Link:        `https://www.phoronix.com/`,
			Description: `Linux Hardware Reviews & News`,
			Language:    "en-us",
			Items: []Item{
				{
					Title:       `AMD Posts Code Enabling "Cyan Skillfish" Display Support Due To Different DCN2 Variant`,
					Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMD-Cyan-Skillfish-DCN-2.01`,
					Description: `Since July we've seen AMD open-source driver engineers posting code for "Cyan Skillfish" as an APU with Navi 1x graphics. While initial support for Cyan Skillfish was merged for Linux 5.15, it turns out the display code isn't yet wired up due to being a different DCN2 variant for its display block...`,
				},
				{
					Title:       `Linux 5.16 To Bring Initial DisplayPort 2.0 Support For AMD Radeon Driver (AMDGPU)`,
					Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMDGPU-DP-2.0-Linux-5.16`,
					Description: `A batch of feature updates was submitted today for DRM-Next of early feature work slated to come to the next version of the Linux kernel...`,
				},
			},
		},
	},
	{
		`<?xml version="1.0" encoding="UTF-8"?><rss version="2.0" xmlns:media="http://search.yahoo.com/mrss/"><channel><title>Noticias en rtve.es</title><description>RSS Tags</description><link>http://www.rtve.es</link><item><title>Erupción en La Palma, en directo | La lava llega a 800 metros del mar y cambia de dirección al norte</title><link>http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml</link><pubDate>Tue, 28 Sep 2021 12:10:00 +0200</pubDate><description>&lt;ul&gt; &lt;li&gt;Varios n&amp;uacute;cleos poblacionales del municipio de Tazacorte han sido confinados&lt;/li&gt; &lt;li&gt;La colada de lava podr&amp;iacute;a llegar a la costa en las pr&amp;oacute;ximas horas&lt;/li&gt; &lt;/ul&gt;&lt;br/&gt;&lt;a href="http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml"&gt;Leer la noticia completa&lt;/a&gt;&lt;img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&amp;cg=F-N-B-TENOTICI-TESESPE01-TES800089&amp;si=http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml" alt=""/&gt;</description></item><item><title>Guía de restricciones COVID: nuevas medidas en ocio nocturno, hostelería y aforos</title><link>http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml</link><pubDate>Tue, 28 Sep 2021 12:02:00 +0200</pubDate><description>&lt;ul&gt; &lt;li&gt;Repasa las principales medidas y restricciones frente a la COVID-19, comunidad a comunidad&lt;/li&gt; &lt;li&gt;&lt;a href="https://www.rtve.es/noticias/20210928/coronavirus-covid-directo-espana-mundo-ultima-hora/2175601.shtml" target="_blank"&gt;Coronavirus: &amp;uacute;ltima hora&lt;/a&gt;&amp;nbsp;|&amp;nbsp;&lt;a href="https://www.rtve.es/noticias/20210924/mapa-del-coronavirus-espana/2004681.shtml" target="_blank"&gt;Mapa de Espa&amp;ntilde;a&lt;/a&gt;&amp;nbsp;|&amp;nbsp;&lt;a href="https://www.rtve.es/noticias/20210924/ocupacion-camas-covid-19-hospitales-espanoles/2042349.shtml" target="_blank"&gt;Hospitales y UCI&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;a href="https://www.rtve.es/noticias/20210924/campana-vacunacion-espana/2062499.shtml" target="_blank"&gt;Vacunas en Espa&amp;ntilde;a&lt;/a&gt;&amp;nbsp;|&amp;nbsp;&lt;a href="https://www.rtve.es/noticias/20210924/mapa-mundial-del-coronavirus/1998143.shtml" target="_blank"&gt;Mapa mundial&amp;#8203;&lt;/a&gt;&amp;nbsp;|&amp;nbsp;&lt;a href="https://www.rtve.es/lab/vacunacion-espana-coronavirus/"&gt;Especial: La gran vacunaci&amp;oacute;n&lt;/a&gt;&lt;/li&gt; &lt;/ul&gt;&lt;br/&gt;&lt;a href="http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml"&gt;Leer la noticia completa&lt;/a&gt;&lt;img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&amp;cg=F-N-B-TENOTICI-TESESPE01-TELCO20VX&amp;si=http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml" alt=""/&gt;</description></item></channel></rss>`,
		Channel{
			Title:       `Noticias en rtve.es`,
			Link:        `http://www.rtve.es`,
			Description: `RSS Tags`,
			Items: []Item{
				{
					Title:       `Erupción en La Palma, en directo | La lava llega a 800 metros del mar y cambia de dirección al norte`,
					Link:        `http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml`,
					Description: `<ul> <li>Varios n&uacute;cleos poblacionales del municipio de Tazacorte han sido confinados</li> <li>La colada de lava podr&iacute;a llegar a la costa en las pr&oacute;ximas horas</li> </ul><br/><a href="http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml">Leer la noticia completa</a><img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&cg=F-N-B-TENOTICI-TESESPE01-TES800089&si=http://www.rtve.es/noticias/20210928/erupcion-palma-directo-lava-llega-800-metros-del-mar-cambia-direccion-norte/2175602.shtml" alt=""/>`,
				},
				{
					Title:       `Guía de restricciones COVID: nuevas medidas en ocio nocturno, hostelería y aforos`,
					Link:        `http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml`,
					Description: `<ul> <li>Repasa las principales medidas y restricciones frente a la COVID-19, comunidad a comunidad</li> <li><a href="https://www.rtve.es/noticias/20210928/coronavirus-covid-directo-espana-mundo-ultima-hora/2175601.shtml" target="_blank">Coronavirus: &uacute;ltima hora</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/mapa-del-coronavirus-espana/2004681.shtml" target="_blank">Mapa de Espa&ntilde;a</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/ocupacion-camas-covid-19-hospitales-espanoles/2042349.shtml" target="_blank">Hospitales y UCI</a></li> <li><a href="https://www.rtve.es/noticias/20210924/campana-vacunacion-espana/2062499.shtml" target="_blank">Vacunas en Espa&ntilde;a</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/noticias/20210924/mapa-mundial-del-coronavirus/1998143.shtml" target="_blank">Mapa mundial&#8203;</a>&nbsp;|&nbsp;<a href="https://www.rtve.es/lab/vacunacion-espana-coronavirus/">Especial: La gran vacunaci&oacute;n</a></li> </ul><br/><a href="http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml">Leer la noticia completa</a><img src="http://secure-uk.imrworldwide.com/cgi-bin/m?ci=es-rssrtve&cg=F-N-B-TENOTICI-TESESPE01-TELCO20VX&si=http://www.rtve.es/noticias/20210928/guia-restricciones-covid-nuevas-medidas-ocio-nocturno-hosteleria-aforos/2041269.shtml" alt=""/>`,
				},
			},
		},
	},
}
